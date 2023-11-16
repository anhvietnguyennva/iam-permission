package dto

import "iam-permission/internal/pkg/entity"

type CreateSubjectSetRelationTupleRes struct {
	ID string `json:"id"`
}

func (r *CreateSubjectSetRelationTupleRes) FromEntity(tuple *entity.SubjectSetRelationTuple) *CreateSubjectSetRelationTupleRes {
	r.ID = tuple.ID.String()
	return r
}
