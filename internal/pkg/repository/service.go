package repository

import (
	"context"

	"github.com/anhvietnguyennva/go-error/pkg/errors"

	"iam-permission/internal/pkg/entity"
)

type IServiceRepository interface {
	Create(ctx context.Context, service *entity.Service) (*entity.Service, *errors.InfraError)
}
