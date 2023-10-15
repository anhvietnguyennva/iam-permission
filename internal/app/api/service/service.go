package service

import (
	"context"

	"github.com/anhvietnguyennva/go-error/pkg/errors"

	"iam-permission/internal/pkg/entity"
)

type IServiceService interface {
	Create(ctx context.Context, service *entity.Service) (*entity.Service, *errors.DomainError)
}
