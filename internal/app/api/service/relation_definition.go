package service

import (
	"context"

	"github.com/anhvietnguyennva/go-error/pkg/errors"

	"iam-permission/internal/app/api/valueobject"
	"iam-permission/internal/pkg/entity"
)

type IRelationDefinitionService interface {
	Create(ctx context.Context, req *valueobject.CreateRelationDefinitionReq) (*entity.RelationDefinition, *errors.DomainError)
}
