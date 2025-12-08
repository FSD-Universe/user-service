package main

import (
	"auth-service/src/interfaces/content"
	"auth-service/src/server"
	"context"
	"fmt"
	"net"
	"time"

	grpcImpl "auth-service/src/grpc"
	c "auth-service/src/interfaces/config"
	g "auth-service/src/interfaces/global"
	pb "auth-service/src/interfaces/grpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"half-nothing.cn/service-core/cleaner"
	"half-nothing.cn/service-core/config"
	"half-nothing.cn/service-core/discovery"
	"half-nothing.cn/service-core/interfaces/global"
	"half-nothing.cn/service-core/logger"
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
	lg.Init(
		global.LogName,
		applicationConfig.GlobalConfig.LogConfig,
	)

	lg.Info(" _____     _   _   _____             _")
	lg.Info("|  _  |_ _| |_| |_|   __|___ ___ _ _|_|___ ___")
	lg.Info("|     | | |  _|   |__   | -_|  _| | | |  _| -_|")
	lg.Info("|__|__|___|_| |_|_|_____|___|_|  \\_/|_|___|___|")
	lg.Infof("                     Copyright © %d-%d Half_nothing", global.BeginYear, time.Now().Year())
	lg.Infof("                                   AuthService v%s", g.AppVersion)

	cl := cleaner.NewCleaner(lg)
	cl.Init()

	applicationContent := content.NewApplicationContentBuilder().
		SetConfigManager(configManager).
		SetCleaner(cl).
		SetLogger(lg).
		Build()

	go server.StartServer(applicationContent)

	if applicationConfig.ServerConfig.GrpcServerConfig.Enable {
		started := make(chan bool)
		go func() {
			address := fmt.Sprintf("%s:%d", applicationConfig.ServerConfig.GrpcServerConfig.Host, applicationConfig.ServerConfig.GrpcServerConfig.Port)
			lis, err := net.Listen("tcp", address)
			if err != nil {
				lg.Fatalf("gRPC fail to listen: %v", err)
				started <- false
				return
			}
			s := grpc.NewServer()
			grpcServer := grpcImpl.NewAuthServer(lg)
			pb.RegisterAuthServer(s, grpcServer)
			reflection.Register(s)
			cl.Add("gRPC Server", func(ctx context.Context) error {
				timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
				defer cancel()
				cleanOver := make(chan struct{})
				go func() {
					s.GracefulStop()
					cleanOver <- struct{}{}
				}()
				select {
				case <-timeoutCtx.Done():
					s.Stop()
				case <-cleanOver:
				}
				return nil
			})
			lg.Infof("gRPC server listening at %v", lis.Addr())
			started <- true
			if err := s.Serve(lis); err != nil {
				lg.Fatalf("gRPC failed to serve: %v", err)
				return
			}
		}()

		go func() {
			start := <-started
			if !start {
				return
			}
			version := utils.NewVersion(g.AppVersion)
			service := discovery.NewServiceDiscovery(
				lg,
				"user-service",
				applicationConfig.ServerConfig.GrpcServerConfig.Port,
				version,
			)
			if err := service.Start(); err != nil {
				lg.Fatalf("fail to start service discovery: %v", err)
				cl.Clean()
				return
			}
			cl.Add("Service Discovery", service.Stop)
		}()
	}

	cl.Wait()
}
