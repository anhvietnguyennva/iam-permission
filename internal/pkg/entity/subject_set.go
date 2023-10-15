package entity

import "github.com/google/uuid"

type SubjectSet struct {
	BaseID
	BaseCreatedUpdated
	ServiceID            uuid.UUID
	RelationDefinitionID uuid.UUID
	Namespace            string
	Object               string
	Relation             string
	Description          string
}
