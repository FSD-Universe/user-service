package main

import (
	"context"
	"fmt"
	"time"
	"user-service/src/interfaces/content"
	g "user-service/src/interfaces/global"
	"user-service/src/repository"
	"user-service/src/server"

	c "user-service/src/interfaces/config"
	pb "user-service/src/interfaces/grpc"

	capi "github.com/hashicorp/consul/api"
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

	utils.CheckStringEnv(g.EnvEmailServiceName, g.EmailServiceName)
	utils.CheckStringEnv(g.EnvAuditServiceName, g.AuditServiceName)
	utils.CheckIntEnv(g.EnvBcryptCost, g.BcryptCost)

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

	contentBuilder := content.NewApplicationContentBuilder().
		SetConfigManager(configManager).
		SetCleaner(cl).
		SetLogger(lg).
		SetJwtClaimFactory(jwt.NewClaimFactory(applicationConfig.JwtConfig)).
		SetUserRepo(repository.NewUserRepository(lg, db, applicationConfig.DatabaseConfig.QueryTimeoutDuration)).
		SetRoleRepo(repository.NewRoleRepository(lg, db, applicationConfig.DatabaseConfig.QueryTimeoutDuration))

	requiredServices := []string{*g.EmailServiceName, *g.AuditServiceName}

	consulClient := discovery.NewConsulClient(lg, applicationConfig.GlobalConfig.Discovery, g.AppVersion)

	// since we register our service,
	// but actually we do not provide any grpc interface,
	// just for health check
	if err := consulClient.RegisterServer(); err != nil {
		lg.Fatalf("fail to register server: %v", err)
		return
	}

	cl.Add("Discovery", consulClient.UnregisterServer)

	consulClient.WatchServices(requiredServices)

	cl.Add("ServiceWatcher", consulClient.StopWatch)

	if err := consulClient.WaitForServices(*global.ReconnectTimeout); err != nil {
		lg.Fatalf("fail to wait for required services: %v", err)
		return
	}

	lg.Info("all required services are online")

	connManager := grpcUtils.NewClientConnections(lg)
	cl.Add("GrpcClient", connManager.Close)

	clientManager := content.NewGrpcClientManager(nil, nil)
	contentBuilder.SetGrpcClientManager(clientManager)

	listener := discovery.NewServiceListener(
		consulClient.EventChan,
		discovery.KeepRequiredServiceOnline(
			lg,
			consulClient,
			cl.Clean,
			func(serviceName string, info *capi.ServiceEntry) {
				if serviceName == *g.EmailServiceName {
					emailConn, err := grpcUtils.InitGrpcClient(lg, applicationConfig.TelemetryConfig, applicationConfig.ClientConfig, info)
					if err != nil {
						lg.Fatalf("fail to start email grpc client: %v", err)
						cl.Clean()
						return
					}
					connManager.Add(*g.EmailServiceName, emailConn)
					clientManager.SetEmailClient(pb.NewEmailClient(emailConn))
				}
				if serviceName == *g.AuditServiceName {
					auditConn, err := grpcUtils.InitGrpcClient(lg, applicationConfig.TelemetryConfig, applicationConfig.ClientConfig, info)
					if err != nil {
						lg.Fatalf("fail to start audit log grpc client: %v", err)
						cl.Clean()
						return
					}
					connManager.Add(*g.AuditServiceName, auditConn)
					clientManager.SetAuditLogClient(pb.NewAuditLogClient(auditConn))
				}
			},
		),
	)
	listener.Start(context.Background())
	cl.Add("ServiceListener", listener.Stop)

	go server.StartServer(contentBuilder.Build())

	cl.Wait()
}
