package service

import (
	"context"

	"github.com/anhvietnguyennva/go-error/pkg/errors"

	"iam-permission/internal/app/api/valueobject"
	"iam-permission/internal/pkg/entity"
)

type ISubjectSetRelationTupleService interface {
	Create(ctx context.Context, req *valueobject.CreateSubjectSetRelationTupleReq) (*entity.SubjectSetRelationTuple, *errors.DomainError)
}
