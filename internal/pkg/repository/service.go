package repository

import (
	"context"

	"github.com/anhvietnguyennva/go-error/pkg/errors"

	"iam-permission/internal/pkg/entity"
)

type IServiceRepository interface {
	GetByNamespace(ctx context.Context, namespace string) (*entity.Service, *errors.InfraError)

	Create(ctx context.Context, service *entity.Service) (*entity.Service, *errors.InfraError)
}
