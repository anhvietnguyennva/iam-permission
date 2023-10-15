package postgres

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"iam-permission/internal/pkg/entity"
)

type BaseID struct {
	ID uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (b *BaseID) BeforeCreate(_ *gorm.DB) error {
	if b.ID != uuid.Nil {
		return nil
	}
	b.ID = uuid.New()
	return nil
}

func (b *BaseID) ToEntity() entity.BaseID {
	return entity.BaseID{
		ID: b.ID,
	}
}

type BaseCreatedUpdated struct {
	CreatedBy string
	UpdatedBy string
	CreatedAt uint64 `gorm:"autoCreateTime"`
	UpdatedAt uint64 `gorm:"autoUpdateTime"`
}

func (b BaseCreatedUpdated) ToEntity() entity.BaseCreatedUpdated {
	return entity.BaseCreatedUpdated{
		CreatedBy: b.CreatedBy,
		UpdatedBy: b.UpdatedBy,
		CreatedAt: b.CreatedAt,
		UpdatedAt: b.UpdatedAt,
	}
}

func (b BaseCreatedUpdated) FromEntity(e entity.BaseCreatedUpdated) BaseCreatedUpdated {
	b.CreatedBy = e.CreatedBy
	b.UpdatedBy = e.UpdatedBy
	b.CreatedAt = e.CreatedAt
	b.UpdatedAt = e.UpdatedAt
	return b
}
