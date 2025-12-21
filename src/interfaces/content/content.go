// Package content
package content

import (
	c "user-service/src/interfaces/config"
	"user-service/src/interfaces/grpc"
	"user-service/src/interfaces/repository"

	"half-nothing.cn/service-core/interfaces/cleaner"
	"half-nothing.cn/service-core/interfaces/config"
	"half-nothing.cn/service-core/interfaces/http/jwt"
	"half-nothing.cn/service-core/interfaces/logger"
)

// ApplicationContent 应用程序上下文结构体，包含所有核心组件的接口
type ApplicationContent struct {
	configManager  config.ManagerInterface[*c.Config] // 配置管理器
	cleaner        cleaner.Interface                  // 清理器
	logger         logger.Interface                   // 日志
	claimFactory   jwt.ClaimFactoryInterface          // JWT 令牌工厂
	userRepo       repository.UserInterface           // 用户数据库操作
	emailClient    grpc.EmailClient                   // 邮件服务
	auditLogClient grpc.AuditLogClient                // 审计日志服务
}

func (app *ApplicationContent) ConfigManager() config.ManagerInterface[*c.Config] {
	return app.configManager
}

func (app *ApplicationContent) Cleaner() cleaner.Interface { return app.cleaner }

func (app *ApplicationContent) Logger() logger.Interface { return app.logger }

func (app *ApplicationContent) ClaimFactory() jwt.ClaimFactoryInterface {
	return app.claimFactory
}

func (app *ApplicationContent) UserRepo() repository.UserInterface {
	return app.userRepo
}

func (app *ApplicationContent) EmailClient() grpc.EmailClient {
	return app.emailClient
}

func (app *ApplicationContent) AuditLogClient() grpc.AuditLogClient {
	return app.auditLogClient
}
