package dto

import (
	"strings"

	"iam-permission/internal/app/api/valueobject"
)

type CreateRelationDefinitionReq struct {
	Namespace   string `json:"namespace" binding:"required,max=255"`
	Relation    string `json:"relation" binding:"required,max=255"`
	Description string `json:"description" binding:"max=255"`
}

func (r *CreateRelationDefinitionReq) ToValueObject() *valueobject.CreateRelationDefinitionReq {
	return &valueobject.CreateRelationDefinitionReq{
		Namespace:   strings.ToLower(r.Namespace),
		Relation:    strings.ToLower(r.Relation),
		Description: r.Description,
	}
}
