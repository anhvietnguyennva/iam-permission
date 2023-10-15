package postgres

import (
	"github.com/google/uuid"

	"iam-permission/internal/pkg/entity"
)

type RelationConfiguration struct {
	BaseID
	BaseCreatedUpdated

	ServiceID                  uuid.UUID `gorm:"column:service_id;<-:create"`
	Namespace                  string    `gorm:"<-:create"`
	ParentRelationDefinitionID uuid.UUID `gorm:"column:parent_relation_definition_id;<-:create"`
	ParentRelation             string    `gorm:"<-:create"`
	ChildRelationDefinitionID  uuid.UUID `gorm:"column:child_relation_definition_id;<-:create"`
	ChildRelation              string    `gorm:"<-:create"`
}

func (c *RelationConfiguration) TableName() string {
	return "relation_configurations"
}

func (c *RelationConfiguration) FromEntity(e *entity.RelationConfiguration) *RelationConfiguration {
	c.ID = e.ID
	c.ServiceID = e.ServiceID
	c.Namespace = e.Namespace
	c.ParentRelationDefinitionID = e.ParentRelationDefinitionID
	c.ParentRelation = e.ParentRelation
	c.ChildRelationDefinitionID = e.ChildRelationDefinitionID
	c.ChildRelation = e.ChildRelation
	c.BaseCreatedUpdated = BaseCreatedUpdated{}.FromEntity(e.BaseCreatedUpdated)
	return c
}

func (c *RelationConfiguration) ToEntity() *entity.RelationConfiguration {
	return &entity.RelationConfiguration{
		BaseID:                     c.BaseID.ToEntity(),
		BaseCreatedUpdated:         c.BaseCreatedUpdated.ToEntity(),
		ServiceID:                  c.ServiceID,
		Namespace:                  c.Namespace,
		ParentRelationDefinitionID: c.ParentRelationDefinitionID,
		ParentRelation:             c.ParentRelation,
		ChildRelationDefinitionID:  c.ChildRelationDefinitionID,
		ChildRelation:              c.ChildRelation,
	}
}
