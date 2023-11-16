package dto

import (
	"strings"

	"iam-permission/internal/app/api/valueobject"
)

type CheckPermissionReq struct {
	Namespace string `form:"namespace" binding:"required,max=255"`
	Object    string `form:"object" binding:"required,max=255"`
	Relation  string `form:"relation" binding:"required,max=255"`
	SubjectID string `form:"subjectId" binding:"required,max=255"`
	MaxDepth  uint8  `form:"maxDepth,default=3" binding:"max=3"`
}

func (r *CheckPermissionReq) ToValueObject() *valueobject.CheckPermissionReq {
	return &valueobject.CheckPermissionReq{
		Namespace: strings.ToLower(r.Namespace),
		Object:    strings.ToLower(r.Object),
		Relation:  strings.ToLower(r.Relation),
		SubjectID: strings.ToLower(r.SubjectID),
		MaxDepth:  r.MaxDepth,
	}
}
