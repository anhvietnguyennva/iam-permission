package repository

import (
	"context"

	"github.com/anhvietnguyennva/go-error/pkg/errors"

	"iam-permission/internal/pkg/entity"
)

type ISubjectRelationTupleRepository interface {
	Create(ctx context.Context, tuple *entity.SubjectRelationTuple) (*entity.SubjectRelationTuple, *errors.InfraError)

	GetByNamespaceAndObjectAndRelationAndSubjectID(ctx context.Context, namespace string, object string, relation string, subjectID string) (*entity.SubjectRelationTuple, *errors.InfraError)
}
