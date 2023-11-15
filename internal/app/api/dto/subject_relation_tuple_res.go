package dto

import "iam-permission/internal/pkg/entity"

type CreateSubjectRelationTupleRes struct {
	ID string `json:"id"`
}

func (r *CreateSubjectRelationTupleRes) FromEntity(tuple *entity.SubjectRelationTuple) *CreateSubjectRelationTupleRes {
	r.ID = tuple.ID.String()
	return r
}
