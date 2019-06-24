package models

import (
	"time"

	"github.com/google/uuid"
)

type Base struct {
	ID        string     `gorm:"type:varchar(36);primary_key" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" sql:"index"`
}

// generate uuid using gorm hooks
func (m *Base) BeforeCreate() {
	u, _ := uuid.NewUUID()
	m.ID = u.String()
}
