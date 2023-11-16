package dto

import (
	"strings"

	"iam-permission/internal/app/api/valueobject"
)

type CreateSubjectSetRelationTupleReq struct {
	Namespace           string `json:"namespace" binding:"required,max=255"`
	Object              string `json:"object" binding:"required,max=255"`
	Relation            string `json:"relation" binding:"required,max=255"`
	SubjectSetNamespace string `json:"subjectSetNamespace" binding:"required,max=255"`
	SubjectSetObject    string `json:"subjectSetObject" binding:"required,max=255"`
	SubjectSetRelation  string `json:"SubjectSetRelation" binding:"required,max=255"`
}

func (r *CreateSubjectSetRelationTupleReq) ToValueObject() *valueobject.CreateSubjectSetRelationTupleReq {
	return &valueobject.CreateSubjectSetRelationTupleReq{
		Namespace:           strings.ToLower(r.Namespace),
		Object:              strings.ToLower(r.Object),
		Relation:            strings.ToLower(r.Relation),
		SubjectSetNamespace: strings.ToLower(r.SubjectSetNamespace),
		SubjectSetObject:    strings.ToLower(r.SubjectSetObject),
		SubjectSetRelation:  strings.ToLower(r.SubjectSetRelation),
	}
}
