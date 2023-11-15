package repository

import (
	"context"

	"github.com/anhvietnguyennva/go-error/pkg/errors"

	"iam-permission/internal/pkg/entity"
)

type ISubjectRelationTupleRepository interface {
	Create(ctx context.Context, tuple *entity.SubjectRelationTuple) (*entity.SubjectRelationTuple, *errors.InfraError)
}
