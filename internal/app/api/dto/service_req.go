package dto

type CreateServiceRequest struct {
	Namespace   string `json:"namespace" binding:"required,max=255"`
	Description string `json:"description" binding:"required,max=255"`
}
