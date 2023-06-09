package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Base struct {
	ID        uuid.UUID      `gorm:"type:char(36);primary-key" json:"id,omitempty"`
	CreatedAt time.Time      `json:"created_at,omitempty"`
	UpdatedAt time.Time      `json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (base *Base) EmptyID() bool {
	return base.ID == uuid.UUID{}
}

func (base *Base) BeforeCreate(tx *gorm.DB) error {
	if base.EmptyID() {
		base.ID = uuid.New()
	}
	return nil
}
