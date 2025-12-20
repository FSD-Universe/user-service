// Package content
package content

import (
	c "user-service/src/interfaces/config"
	"user-service/src/interfaces/repository"

	"half-nothing.cn/service-core/interfaces/cleaner"
	"half-nothing.cn/service-core/interfaces/config"
	"half-nothing.cn/service-core/interfaces/http/jwt"
	"half-nothing.cn/service-core/interfaces/logger"
)

type ApplicationContentBuilder struct {
	content *ApplicationContent
}

func NewApplicationContentBuilder() *ApplicationContentBuilder {
	return &ApplicationContentBuilder{
		content: &ApplicationContent{},
	}
}

func (builder *ApplicationContentBuilder) SetConfigManager(configManager config.ManagerInterface[*c.Config]) *ApplicationContentBuilder {
	builder.content.configManager = configManager
	return builder
}

func (builder *ApplicationContentBuilder) SetCleaner(cleaner cleaner.Interface) *ApplicationContentBuilder {
	builder.content.cleaner = cleaner
	return builder
}

func (builder *ApplicationContentBuilder) SetLogger(logger logger.Interface) *ApplicationContentBuilder {
	builder.content.logger = logger
	return builder
}

func (builder *ApplicationContentBuilder) SetUserRepo(userRepo repository.UserInterface) *ApplicationContentBuilder {
	builder.content.userRepo = userRepo
	return builder
}

func (builder *ApplicationContentBuilder) SetJwtClaimFactory(claimFactory jwt.ClaimFactoryInterface) *ApplicationContentBuilder {
	builder.content.claimFactory = claimFactory
	return builder
}

func (builder *ApplicationContentBuilder) Build() *ApplicationContent {
	return builder.content
}
