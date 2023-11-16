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

type RelationDefinitionService struct {
	serviceRepo            repo.IServiceRepository
	relationDefinitionRepo repo.IRelationDefinitionRepository
}

var relationDefinitionService *RelationDefinitionService

func InitRelationDefinitionService(
	serviceRepo repo.IServiceRepository,
	relationDefinitionRepo repo.IRelationDefinitionRepository,
) {
	lock.InitComponentLock.Lock()
	defer lock.InitComponentLock.Unlock()
	if relationDefinitionService == nil {
		relationDefinitionService = &RelationDefinitionService{
			serviceRepo:            serviceRepo,
			relationDefinitionRepo: relationDefinitionRepo,
		}
	}
}

func RelationDefinitionServiceInstance() *RelationDefinitionService {
	return relationDefinitionService
}

func (s *RelationDefinitionService) Create(ctx context.Context, req *valueobject.CreateRelationDefinitionReq) (*entity.RelationDefinition, *errors.DomainError) {
	// get service
	service, infraErr := s.serviceRepo.GetByNamespace(ctx, req.Namespace)
	if infraErr != nil {
		return nil, errt.DomainTransformerInstance().InfraErrToDomainErr(infraErr)
	}

	// create
	definition := &entity.RelationDefinition{
		ServiceID:   service.ID,
		Namespace:   service.Namespace,
		Relation:    req.Relation,
		Description: req.Description,
		BaseCreatedUpdated: entity.BaseCreatedUpdated{
			CreatedBy: req.CreatedBy,
			UpdatedBy: req.CreatedBy,
		},
	}
	definition, infraErr = s.relationDefinitionRepo.Create(ctx, definition)
	if infraErr != nil {
		return nil, errt.DomainTransformerInstance().InfraErrToDomainErr(infraErr)
	}

	return definition, nil
}
