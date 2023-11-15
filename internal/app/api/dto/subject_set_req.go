package dto

import (
	"strings"

	"iam-permission/internal/app/api/valueobject"
)

type CreateSubjectSetReq struct {
	Namespace   string `json:"namespace" binding:"required,max=255"`
	Object      string `json:"object" binding:"required,max=255"`
	Relation    string `json:"relation" binding:"required,max=255"`
	Description string `json:"description" binding:"max=500"`
}

func (r *CreateSubjectSetReq) ToValueObject() *valueobject.CreateSubjectSetReq {
	return &valueobject.CreateSubjectSetReq{
		Namespace:   strings.ToLower(r.Namespace),
		Object:      strings.ToLower(r.Object),
		Relation:    strings.ToLower(r.Relation),
		Description: r.Description,
	}
}
