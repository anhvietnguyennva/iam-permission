package entity

import "github.com/google/uuid"

type SubjectRelationTuple struct {
	BaseID
	BaseCreatedUpdated
	ServiceID            uuid.UUID
	RelationDefinitionID uuid.UUID
	Namespace            string
	Object               string
	Relation             string
	SubjectID            string
}
