package service

import (
	"context"

	"github.com/anhvietnguyennva/go-error/pkg/errors"

	"iam-permission/internal/app/api/valueobject"
)

type IPermissionService interface {
	CheckPermission(ctx context.Context, req *valueobject.CheckPermissionReq) (bool, *errors.DomainError)
}
