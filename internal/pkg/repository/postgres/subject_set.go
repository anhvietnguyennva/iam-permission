package postgres

import (
	"github.com/google/uuid"
	"iam-permission/internal/pkg/entity"
)

type SubjectSet struct {
	BaseID
	BaseCreatedUpdated

	ServiceID            uuid.UUID `gorm:"column:service_id;<-:create"`
	RelationDefinitionID uuid.UUID `gorm:"column:relation_definition_id;<-:create"`
	Namespace            string    `gorm:"<-:create"`
	Object               string    `gorm:"<-:create"`
	Relation             string    `gorm:"<-:create"`
	Description          string
}

func (s *SubjectSet) TableName() string {
	return "subject_sets"
}

func (s *SubjectSet) FromEntity(e *entity.SubjectSet) *SubjectSet {
	s.ID = e.ID
	s.CreatedAt = e.CreatedAt
	s.CreatedBy = e.CreatedBy
	s.UpdatedAt = e.UpdatedAt
	s.UpdatedBy = e.UpdatedBy
	s.ServiceID = e.ServiceID
	s.Namespace = e.Namespace
	s.RelationDefinitionID = e.RelationDefinitionID
	s.Object = e.Object
	s.Relation = e.Relation
	s.Description = e.Description
	return s
}

func (s *SubjectSet) ToEntity() *entity.SubjectSet {
	return &entity.SubjectSet{
		BaseID:               s.BaseID.ToEntity(),
		BaseCreatedUpdated:   s.BaseCreatedUpdated.ToEntity(),
		ServiceID:            s.ServiceID,
		Namespace:            s.Namespace,
		RelationDefinitionID: s.RelationDefinitionID,
		Object:               s.Object,
		Relation:             s.Relation,
		Description:          s.Description,
	}
}
