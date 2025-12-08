// Package content
package content

import (
	c "auth-service/src/interfaces/config"

	"half-nothing.cn/service-core/interfaces/cleaner"
	"half-nothing.cn/service-core/interfaces/config"
	"half-nothing.cn/service-core/interfaces/logger"
)

// ApplicationContent 应用程序上下文结构体，包含所有核心组件的接口
type ApplicationContent struct {
	configManager config.ManagerInterface[*c.Config] // 配置管理器
	cleaner       cleaner.Interface                  // 清理器
	logger        logger.Interface                   // 日志
}

func (app *ApplicationContent) ConfigManager() config.ManagerInterface[*c.Config] {
	return app.configManager
}

func (app *ApplicationContent) Cleaner() cleaner.Interface { return app.cleaner }

func (app *ApplicationContent) Logger() logger.Interface { return app.logger }
