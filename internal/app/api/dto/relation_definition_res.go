package dto

import "iam-permission/internal/pkg/entity"

type CreateRelationDefinitionResponse struct {
	ID string `json:"id"`
}

func (r *CreateRelationDefinitionResponse) FromEntity(definition *entity.RelationDefinition) *CreateRelationDefinitionResponse {
	r.ID = definition.ID.String()
	return r
}
