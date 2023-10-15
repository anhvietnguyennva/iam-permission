package entity

import "github.com/google/uuid"

type BaseID struct {
	ID uuid.UUID
}

type BaseCreatedUpdated struct {
	CreatedBy string
	UpdatedBy string
	CreatedAt uint64
	UpdatedAt uint64
}
