package models

import (
	"time"

	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/domain/audit"
)

type AuditLogModel struct {
	ID         string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	TenantID   string    `gorm:"type:uuid;index;not null"`
	UserID     string    `gorm:"type:uuid;index"`
	Action     string    `gorm:"not null"`
	Resource   string    `gorm:"not null"`
	Payload    string    `gorm:"type:text"`
	IPAddress  string
	UserAgent  string
	CreatedAt  time.Time `gorm:"index"`
}

func (AuditLogModel) TableName() string {
	return "audit_logs"
}

func (m *AuditLogModel) ToDomain() *audit.AuditLog {
	return &audit.AuditLog{
		ID:        m.ID,
		TenantID:  m.TenantID,
		UserID:    m.UserID,
		Action:    m.Action,
		Resource:  m.Resource,
		Payload:   m.Payload,
		IPAddress: m.IPAddress,
		UserAgent: m.UserAgent,
		CreatedAt: m.CreatedAt,
	}
}

func FromAuditDomain(l *audit.AuditLog) *AuditLogModel {
	return &AuditLogModel{
		ID:        l.ID,
		TenantID:  l.TenantID,
		UserID:    l.UserID,
		Action:    l.Action,
		Resource:  l.Resource,
		Payload:   l.Payload,
		IPAddress: l.IPAddress,
		UserAgent: l.UserAgent,
		CreatedAt: l.CreatedAt,
	}
}
