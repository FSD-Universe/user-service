package main

import (
	"fmt"
	"time"
	"user-service/src/interfaces/content"
	g "user-service/src/interfaces/global"
	"user-service/src/repository"
	"user-service/src/server"

	grpcImpl "user-service/src/grpc"
	c "user-service/src/interfaces/config"
	pb "user-service/src/interfaces/grpc"

	"google.golang.org/grpc"
	"half-nothing.cn/service-core/cleaner"
	"half-nothing.cn/service-core/config"
	"half-nothing.cn/service-core/database"
	"half-nothing.cn/service-core/discovery"
	grpcUtils "half-nothing.cn/service-core/grpc"
	"half-nothing.cn/service-core/interfaces/global"
	"half-nothing.cn/service-core/jwt"
	"half-nothing.cn/service-core/logger"
	"half-nothing.cn/service-core/telemetry"
	"half-nothing.cn/service-core/utils"
)

func main() {
	global.CheckFlags()

	configManager := config.NewManager[*c.Config]()
	if err := configManager.Init(); err != nil {
		fmt.Printf("fail to initialize configuration file: %v", err)
		return
	}

	applicationConfig := configManager.GetConfig()
	lg := logger.NewLogger()
	lg.Init(global.LogName, applicationConfig.GlobalConfig.LogConfig)

	lg.Info(" _____             _____             _")
	lg.Info("|  |  |___ ___ ___|   __|___ ___ _ _|_|___ ___")
	lg.Info("|  |  |_ -| -_|  _|__   | -_|  _| | | |  _| -_|")
	lg.Info("|_____|___|___|_| |_____|___|_|  \\_/|_|___|___|")
	lg.Info(fmt.Sprintf("%47s", fmt.Sprintf("Copyright © %d-%d Half_nothing", global.BeginYear, time.Now().Year())))
	lg.Info(fmt.Sprintf("%47s", fmt.Sprintf("UserService v%s", g.AppVersion)))

	cl := cleaner.NewCleaner(lg)
	cl.Init()

	defer cl.Clean()

	closeFunc, db, err := database.InitDatabase(lg, applicationConfig.DatabaseConfig)
	if err != nil {
		lg.Fatalf("fail to initialize database: %v", err)
		return
	}
	cl.Add("Database", closeFunc)

	if applicationConfig.TelemetryConfig.Enable {
		if err := telemetry.InitSDK(lg, cl, applicationConfig.TelemetryConfig); err != nil {
			lg.Fatalf("fail to initialize telemetry: %v", err)
			return
		}

		if applicationConfig.TelemetryConfig.DatabaseTrace {
			if err := database.ApplyDBTracing(db, "mysql"); err != nil {
				lg.Fatalf("fail to apply database tracing: %v", err)
				return
			}
		}
	}

	applicationContent := content.NewApplicationContentBuilder().
		SetConfigManager(configManager).
		SetCleaner(cl).
		SetLogger(lg).
		SetJwtClaimFactory(jwt.NewClaimFactory(applicationConfig.JwtConfig)).
		SetUserRepo(repository.NewUserRepository(lg, db, applicationConfig.DatabaseConfig.QueryTimeoutDuration)).
		Build()

	go server.StartServer(applicationContent)

	if applicationConfig.ServerConfig.GrpcServerConfig.Enable {
		started := make(chan bool)
		initFunc := func(s *grpc.Server) {
			grpcServer := grpcImpl.NewAuthServer(applicationContent)
			pb.RegisterAuthServer(s, grpcServer)
		}
		if applicationConfig.TelemetryConfig.Enable && applicationConfig.TelemetryConfig.GrpcServerTrace {
			go grpcUtils.StartGrpcServerWithTrace(lg, cl, applicationConfig.ServerConfig.GrpcServerConfig, started, initFunc)
		} else {
			go grpcUtils.StartGrpcServer(lg, cl, applicationConfig.ServerConfig.GrpcServerConfig, started, initFunc)
		}
		go discovery.StartServiceDiscovery(lg, cl, started, utils.NewVersion(g.AppVersion),
			"user-service", applicationConfig.ServerConfig.GrpcServerConfig.Port)
	}

	cl.Wait()
}
