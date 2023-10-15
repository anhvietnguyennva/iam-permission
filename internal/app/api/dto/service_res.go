package dto

import "iam-permission/internal/pkg/entity"

type CreateServiceResponse struct {
	ID string `json:"id"`
}

func (r *CreateServiceResponse) FromEntity(service *entity.Service) *CreateServiceResponse {
	r.ID = service.ID.String()
	return r
}
