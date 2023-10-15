package postgres

import (
	"github.com/google/uuid"

	"iam-permission/internal/pkg/entity"
)

type SubjectRelationTuple struct {
	BaseID
	BaseCreatedUpdated

	ServiceID            uuid.UUID `gorm:"column:service_id;<-:create"`
	RelationDefinitionID uuid.UUID `gorm:"column:relation_definition_id;<-:create"`
	Namespace            string    `gorm:"<-:create"`
	Object               string    `gorm:"<-:create"`
	Relation             string    `gorm:"<-:create"`
	SubjectID            string    `gorm:"<-:create"`
}

func (t *SubjectRelationTuple) TableName() string {
	return "subject_relation_tuples"
}

func (t *SubjectRelationTuple) FromEntity(e *entity.SubjectRelationTuple) *SubjectRelationTuple {
	t.ID = e.ID
	t.CreatedAt = e.CreatedAt
	t.CreatedBy = e.CreatedBy
	t.UpdatedAt = e.UpdatedAt
	t.UpdatedBy = e.UpdatedBy
	t.ServiceID = e.ServiceID
	t.Namespace = e.Namespace
	t.RelationDefinitionID = e.RelationDefinitionID
	t.Object = e.Object
	t.Relation = e.Relation
	t.SubjectID = e.SubjectID
	return t
}

func (t *SubjectRelationTuple) ToEntity() *entity.SubjectRelationTuple {
	return &entity.SubjectRelationTuple{
		BaseID:               t.BaseID.ToEntity(),
		BaseCreatedUpdated:   t.BaseCreatedUpdated.ToEntity(),
		ServiceID:            t.ServiceID,
		Namespace:            t.Namespace,
		RelationDefinitionID: t.RelationDefinitionID,
		Object:               t.Object,
		Relation:             t.Relation,
		SubjectID:            t.SubjectID,
	}
}
