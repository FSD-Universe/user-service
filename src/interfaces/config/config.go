// Package config
package config

import "half-nothing.cn/service-core/interfaces/config"

type Config struct {
	GlobalConfig *GlobalConfig        `yaml:"global"`
	ServerConfig *config.ServerConfig `yaml:"server"`
	JwtConfig    *config.JwtConfig    `yaml:"jwt"`
}

func (c *Config) InitDefaults() {
	c.GlobalConfig = &GlobalConfig{}
	c.GlobalConfig.InitDefaults()
	c.ServerConfig = &config.ServerConfig{}
	c.ServerConfig.InitDefaults()
	c.JwtConfig = &config.JwtConfig{}
	c.JwtConfig.InitDefaults()
}

func (c *Config) Verify() (bool, error) {
	if ok, err := c.GlobalConfig.Verify(); !ok {
		return ok, err
	}
	if ok, err := c.ServerConfig.Verify(); !ok {
		return ok, err
	}
	return c.JwtConfig.Verify()
}
