package dto

import (
	"strings"

	"iam-permission/internal/app/api/valueobject"
)

type CreateServiceRequest struct {
	Namespace   string `json:"namespace" binding:"required,max=255"`
	Description string `json:"description" binding:"required,max=255"`
}

func (r *CreateServiceRequest) ToValueObject() *valueobject.CreateServiceRequest {
	return &valueobject.CreateServiceRequest{
		Namespace:   strings.ToLower(r.Namespace),
		Description: r.Description,
	}
}
