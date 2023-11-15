package dto

import (
	"strings"

	"iam-permission/internal/app/api/valueobject"
)

type CreateSubjectRelationTupleReq struct {
	Namespace string `json:"namespace" binding:"required,max=255"`
	Object    string `json:"object" binding:"required,max=255"`
	Relation  string `json:"relation" binding:"required,max=255"`
	SubjectID string `json:"subjectId" binding:"required,max=255"`
}

func (r *CreateSubjectRelationTupleReq) ToValueObject() *valueobject.CreateSubjectRelationTupleReq {
	return &valueobject.CreateSubjectRelationTupleReq{
		Namespace: strings.ToLower(r.Namespace),
		Object:    strings.ToLower(r.Object),
		Relation:  strings.ToLower(r.Relation),
		SubjectID: strings.ToLower(r.SubjectID),
	}
}
