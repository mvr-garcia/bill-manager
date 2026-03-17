package persistence

import (
	"time"

	"github.com/google/uuid"
)

// BillModel represents the database model for a bill.
type BillModel struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey"`
	Description string         `gorm:"type:varchar(500);not null"`
	Amount      float64        `gorm:"type:decimal(15,2);not null"`
	DueDate     time.Time      `gorm:"type:date;not null"`
	Status      string         `gorm:"type:varchar(50);not null;default:'pending'"`
	CreatedBy   uuid.UUID      `gorm:"type:uuid;not null"`
	ApprovedBy  *uuid.UUID     `gorm:"type:uuid"`
	CreatedAt   time.Time      `gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`
	Audits      []BillAuditModel `gorm:"foreignKey:BillID;constraint:OnDelete:CASCADE"`
}

// TableName specifies the table name for BillModel.
func (BillModel) TableName() string {
	return "bills"
}

// BillAuditModel represents the database model for a bill audit log.
type BillAuditModel struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey"`
	BillID     uuid.UUID `gorm:"type:uuid;not null;index"`
	Action     string    `gorm:"type:varchar(50);not null;index"`
	PerformedBy uuid.UUID `gorm:"type:uuid;not null;index"`
	IPAddress  string    `gorm:"type:varchar(45)"`
	UserAgent  string    `gorm:"type:varchar(500)"`
	CreatedAt  time.Time `gorm:"autoCreateTime;index"`
	Bill       *BillModel `gorm:"foreignKey:BillID;constraint:OnDelete:CASCADE"`
}

// TableName specifies the table name for BillAuditModel.
func (BillAuditModel) TableName() string {
	return "bill_audits"
}
