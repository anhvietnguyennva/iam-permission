package repository

import (
	"context"

	"github.com/anhvietnguyennva/go-error/pkg/errors"

	"iam-permission/internal/pkg/entity"
)

type ISubjectSetRepository interface {
	Create(ctx context.Context, set *entity.SubjectSet) (*entity.SubjectSet, *errors.InfraError)

	GetByNamespaceAndObjectAndRelation(ctx context.Context, namespace string, object string, relation string) (*entity.SubjectSet, *errors.InfraError)
}
