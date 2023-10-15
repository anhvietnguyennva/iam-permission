package entity

import "github.com/google/uuid"

type RelationDefinition struct {
	BaseID
	BaseCreatedUpdated
	ServiceID   uuid.UUID
	Namespace   string
	Relation    string
	Description string
}
