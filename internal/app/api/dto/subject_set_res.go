package dto

import "iam-permission/internal/pkg/entity"

type CreateSubjectSetRes struct {
	ID string `json:"id"`
}

func (r *CreateSubjectSetRes) FromEntity(s *entity.SubjectSet) *CreateSubjectSetRes {
	r.ID = s.ID.String()
	return r
}
