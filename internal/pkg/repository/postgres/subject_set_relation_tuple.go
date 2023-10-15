package postgres

import (
	"github.com/google/uuid"

	"iam-permission/internal/pkg/entity"
)

type SubjectSetRelationTuple struct {
	BaseID
	BaseCreatedUpdated

	ServiceID            uuid.UUID `gorm:"column:service_id;<-:create"`
	SubjectSetID         uuid.UUID `gorm:"column:subject_set_id;<-:create"`
	RelationDefinitionID uuid.UUID `gorm:"column:relation_definition_id;<-:create"`
	Namespace            string    `gorm:"<-:create"`
	Object               string    `gorm:"<-:create"`
	Relation             string    `gorm:"<-:create"`
	SubjectSetNamespace  string    `gorm:"<-:create"`
	SubjectSetObject     string    `gorm:"<-:create"`
	SubjectSetRelation   string    `gorm:"<-:create"`
}

func (t *SubjectSetRelationTuple) TableName() string {
	return "subject_set_relation_tuples"
}

func (t *SubjectSetRelationTuple) FromEntity(e *entity.SubjectSetRelationTuple) *SubjectSetRelationTuple {
	t.ID = e.ID
	t.CreatedAt = e.CreatedAt
	t.CreatedBy = e.CreatedBy
	t.UpdatedAt = e.UpdatedAt
	t.UpdatedBy = e.UpdatedBy
	t.ServiceID = e.ServiceID
	t.SubjectSetID = e.SubjectSetID
	t.RelationDefinitionID = e.RelationDefinitionID
	t.Namespace = e.Namespace
	t.Object = e.Object
	t.Relation = e.Relation
	t.SubjectSetNamespace = e.SubjectSetNamespace
	t.SubjectSetObject = e.SubjectSetObject
	t.SubjectSetRelation = e.SubjectSetRelation
	return t
}

func (t *SubjectSetRelationTuple) ToEntity() *entity.SubjectSetRelationTuple {
	return &entity.SubjectSetRelationTuple{
		BaseID:               t.BaseID.ToEntity(),
		BaseCreatedUpdated:   t.BaseCreatedUpdated.ToEntity(),
		ServiceID:            t.ServiceID,
		Namespace:            t.Namespace,
		SubjectSetID:         t.SubjectSetID,
		RelationDefinitionID: t.RelationDefinitionID,
		Object:               t.Object,
		Relation:             t.Relation,
		SubjectSetNamespace:  t.SubjectSetNamespace,
		SubjectSetObject:     t.SubjectSetObject,
		SubjectSetRelation:   t.SubjectSetRelation,
	}
}
