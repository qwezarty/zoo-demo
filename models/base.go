package models

import "time"

type Base struct {
	ID        string `gorm:"type:varchar(36);primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}
