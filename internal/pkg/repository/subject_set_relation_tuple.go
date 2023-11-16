package repository

import (
	"context"

	"github.com/anhvietnguyennva/go-error/pkg/errors"

	"iam-permission/internal/pkg/entity"
)

type ISubjectSetRelationTupleRepository interface {
	Create(ctx context.Context, tuple *entity.SubjectSetRelationTuple) (*entity.SubjectSetRelationTuple, *errors.InfraError)

	GetByNamespaceAndObjectAndRelation(ctx context.Context, namespace string, object string, relation string) ([]*entity.SubjectSetRelationTuple, *errors.InfraError)
}
