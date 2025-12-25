// Copyright (c) 2025 Half_nothing
// SPDX-License-Identifier: MIT

// Package content
package content

import "user-service/src/interfaces/grpc"

type GrpcClientManager struct {
	auditLogClient grpc.AuditLogClient
	emailClient    grpc.EmailClient
}

func NewGrpcClientManager(
	auditLogClient grpc.AuditLogClient,
	emailClient grpc.EmailClient,
) *GrpcClientManager {
	return &GrpcClientManager{
		auditLogClient: auditLogClient,
		emailClient:    emailClient,
	}
}

func (manager *GrpcClientManager) SetAuditLogClient(auditLogClient grpc.AuditLogClient) {
	manager.auditLogClient = auditLogClient
}

func (manager *GrpcClientManager) AuditLogClient() grpc.AuditLogClient {
	return manager.auditLogClient
}

func (manager *GrpcClientManager) SetEmailClient(emailClient grpc.EmailClient) {
	manager.emailClient = emailClient
}

func (manager *GrpcClientManager) EmailClient() grpc.EmailClient {
	return manager.emailClient
}
