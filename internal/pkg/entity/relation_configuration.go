package entity

import "github.com/google/uuid"

type RelationConfiguration struct {
	BaseID
	BaseCreatedUpdated
	ServiceID                  uuid.UUID
	Namespace                  string
	ParentRelationDefinitionID uuid.UUID
	ParentRelation             string
	ChildRelationDefinitionID  uuid.UUID
	ChildRelation              string
}
