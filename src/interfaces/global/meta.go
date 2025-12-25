// Package global
package global

import (
	"flag"
	"time"
)

var (
	WaitServiceTimeout = flag.Duration("wait-service-timeout", 5*time.Second, "wait service timeout")
	BcryptCost         = flag.Int("bcrypt-cost", 12, "bcrypt cost")
)

const (
	AppVersion    = "0.1.0"
	ConfigVersion = "0.1.0"

	EmailServiceName    = "email-service"
	AuditLogServiceName = "audit-service"

	EnvWaitServiceTimeout = "WAIT_SERVICE_TIMEOUT"
	EnvBcryptCost         = "BCRYPT_COST"
)
