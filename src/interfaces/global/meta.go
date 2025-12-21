// Package global
package global

import (
	"flag"
	"time"
)

var (
	WaitServiceTimeout = flag.Duration("wait-service-timeout", 5*time.Second, "wait service timeout")
)

const (
	AppVersion    = "0.1.0"
	ConfigVersion = "0.1.0"

	ServiceName         = "user-service"
	EmailServiceName    = "email-service"
	AuditLogServiceName = "audit-service"

	EnvWaitServiceTimeout = "WAIT_SERVICE_TIMEOUT"
)
