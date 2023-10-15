package entity

import "github.com/google/uuid"

type SubjectSetRelationTuple struct {
	BaseID
	BaseCreatedUpdated
	ServiceID            uuid.UUID
	SubjectSetID         uuid.UUID
	RelationDefinitionID uuid.UUID
	Namespace            string
	Object               string
	Relation             string
	SubjectSetNamespace  string
	SubjectSetObject     string
	SubjectSetRelation   string
}
