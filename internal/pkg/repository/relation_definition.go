package repository

import (
	"context"

	"github.com/anhvietnguyennva/go-error/pkg/errors"

	"iam-permission/internal/pkg/entity"
)

type IRelationDefinitionRepository interface {
	GetByNamespaceAndRelation(ctx context.Context, namespace string, relation string) (*entity.RelationDefinition, *errors.InfraError)
}
