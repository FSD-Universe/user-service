// Package config
package config

import "half-nothing.cn/service-core/interfaces/config"

type Config struct {
	GlobalConfig    *GlobalConfig            `yaml:"global"`
	ServerConfig    *config.ServerConfig     `yaml:"server"`
	ClientConfig    *config.GrpcClientConfig `yaml:"client"`
	JwtConfig       *config.JwtConfig        `yaml:"jwt"`
	DatabaseConfig  *config.DatabaseConfig   `yaml:"database"`
	TelemetryConfig *config.TelemetryConfig  `yaml:"telemetry"`
}

func (c *Config) InitDefaults() {
	c.GlobalConfig = &GlobalConfig{}
	c.GlobalConfig.InitDefaults()
	c.ServerConfig = &config.ServerConfig{}
	c.ServerConfig.InitDefaults()
	c.ClientConfig = &config.GrpcClientConfig{}
	c.ClientConfig.InitDefaults()
	c.JwtConfig = &config.JwtConfig{}
	c.JwtConfig.InitDefaults()
	c.DatabaseConfig = &config.DatabaseConfig{}
	c.DatabaseConfig.InitDefaults()
	c.TelemetryConfig = &config.TelemetryConfig{}
	c.TelemetryConfig.InitDefaults()
}

func (c *Config) Verify() (bool, error) {
	if ok, err := c.GlobalConfig.Verify(); !ok {
		return ok, err
	}
	if ok, err := c.ServerConfig.Verify(); !ok {
		return ok, err
	}
	if ok, err := c.ClientConfig.Verify(); !ok {
		return ok, err
	}
	if ok, err := c.JwtConfig.Verify(); !ok {
		return ok, err
	}
	if ok, err := c.DatabaseConfig.Verify(); !ok {
		return ok, err
	}
	if ok, err := c.TelemetryConfig.Verify(); !ok {
		return ok, err
	}
	return c.JwtConfig.Verify()
}
