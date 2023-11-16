package repository

import (
	"context"

	"github.com/anhvietnguyennva/go-error/pkg/errors"
)

type IRelationConfigurationRepository interface {
	GetChildrenRelations(ctx context.Context, namespace string, parentRelation string) ([]string, *errors.InfraError)
}
