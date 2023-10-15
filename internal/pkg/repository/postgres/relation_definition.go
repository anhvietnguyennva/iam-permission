package postgres

import (
	"github.com/google/uuid"
	"iam-permission/internal/pkg/entity"
)

type RelationDefinition struct {
	BaseID
	BaseCreatedUpdated

	ServiceID   uuid.UUID `gorm:"column:service_id;<-:create"`
	Namespace   string    `gorm:"<-:create"`
	Relation    string    `gorm:"<-:create"`
	Description string
}

func (d *RelationDefinition) TableName() string {
	return "relation_definitions"
}

func (d *RelationDefinition) FromEntity(e *entity.RelationDefinition) *RelationDefinition {
	d.ID = e.ID
	d.CreatedAt = e.CreatedAt
	d.CreatedBy = e.CreatedBy
	d.UpdatedAt = e.UpdatedAt
	d.UpdatedBy = e.UpdatedBy
	d.ServiceID = e.ServiceID
	d.Namespace = e.Namespace
	d.Relation = e.Relation
	d.Description = e.Description
	return d
}

func (d *RelationDefinition) ToEntity() *entity.RelationDefinition {
	return &entity.RelationDefinition{
		BaseID:             d.BaseID.ToEntity(),
		BaseCreatedUpdated: d.BaseCreatedUpdated.ToEntity(),
		ServiceID:          d.ServiceID,
		Namespace:          d.Namespace,
		Relation:           d.Relation,
		Description:        d.Description,
	}
}
