package service

import (
	"context"

	"github.com/anhvietnguyennva/go-error/pkg/errors"
	errt "github.com/anhvietnguyennva/go-error/pkg/transformer"

	"iam-permission/internal/app/api/valueobject"
	"iam-permission/internal/pkg/entity"
	"iam-permission/internal/pkg/lock"
	repo "iam-permission/internal/pkg/repository"
)

type ServiceService struct {
	serviceRepo repo.IServiceRepository
}

var serviceSvc *ServiceService

func InitServiceService(serviceRepo repo.IServiceRepository) {
	lock.InitComponentLock.Lock()
	defer lock.InitComponentLock.Unlock()
	if serviceSvc == nil {
		serviceSvc = &ServiceService{
			serviceRepo: serviceRepo,
		}
	}
}

func ServiceServiceInstance() *ServiceService {
	return serviceSvc
}

func (s *ServiceService) Create(ctx context.Context, req *valueobject.CreateServiceRequest) (*entity.Service, *errors.DomainError) {
	service := &entity.Service{
		Namespace:   req.Namespace,
		Description: req.Description,
		BaseCreatedUpdated: entity.BaseCreatedUpdated{
			CreatedBy: req.CreatedBy,
			UpdatedBy: req.CreatedBy,
		},
	}
	service, infraErr := s.serviceRepo.Create(ctx, service)
	if infraErr != nil {
		return nil, errt.DomainTransformerInstance().InfraErrToDomainErr(infraErr)
	}
	return service, nil
}
