// Package global
package global

import (
	"flag"
)

var (
	EmailServiceName = flag.String("email_service_name", "email-service", "email service name")
	AuditServiceName = flag.String("audit_service_name", "audit-service", "audit service name")
	BcryptCost       = flag.Int("bcrypt_cost", 12, "bcrypt cost")
)

const (
	AppVersion    = "0.1.0"
	ConfigVersion = "0.1.0"

	EnvEmailServiceName = "EMAIL_SERVICE_NAME"
	EnvAuditServiceName = "AUDIT_SERVICE_NAME"
	EnvBcryptCost       = "BCRYPT_COST"
)
