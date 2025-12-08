// Package config
package config

import (
	"auth-service/src/interfaces/global"

	"half-nothing.cn/service-core/interfaces/config"
)

type GlobalConfig struct {
	*config.GlobalConfig `yaml:",inline"`
}

func (g *GlobalConfig) InitDefaults() {
	g.GlobalConfig = &config.GlobalConfig{}
	g.GlobalConfig.InitDefaults()
	g.Name = "user-service"
	g.Version = global.ConfigVersion
	g.ConfigVersion = global.ConfigVersion
}
