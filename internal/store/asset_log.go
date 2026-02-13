package store

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type AssetAction string

const (
	ActionCreated  AssetAction = "CREATED"
	ActionUpdated  AssetAction = "UPDATED"
	ActionAssigned AssetAction = "ASSIGNED"
	ActionReturned AssetAction = "RETURNED"
	ActionDeleted  AssetAction = "DELETED"
)

type AssetLog struct {
	ID int64 `gorm:"primaryKey"`

	AssetID int64 `gorm:"index;not null"`
	Asset   Asset `gorm:"constraint:OnDelete:CASCADE;"`

	PerformedByID int64 `gorm:"not null"`
	PerformedBy   User  `gorm:"constraint:OnDelete:SET NULL;"`

	Action  AssetAction `gorm:"type:varchar(30);not null"`
	Details string      `gorm:"type:text"`

	CreatedAt time.Time
}

type AssetLogStore struct {
	db *gorm.DB
}

func (s AssetLogStore) Create(ctx context.Context, asset *AssetLog) error {
	return s.db.WithContext(ctx).Create(asset).Error
}
