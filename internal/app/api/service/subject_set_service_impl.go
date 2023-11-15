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

type SubjectSetService struct {
	serviceRepo            repo.IServiceRepository
	relationDefinitionRepo repo.IRelationDefinitionRepository
	subjectSetRepo         repo.ISubjectSetRepository
}

var subjectSetService *SubjectSetService

func InitSubjectSetService(
	serviceRepo repo.IServiceRepository,
	relationDefinitionRepo repo.IRelationDefinitionRepository,
	subjectSetRepo repo.ISubjectSetRepository,
) {
	lock.InitComponentLock.Lock()
	defer lock.InitComponentLock.Unlock()

	if subjectSetService == nil {
		subjectSetService = &SubjectSetService{
			serviceRepo:            serviceRepo,
			relationDefinitionRepo: relationDefinitionRepo,
			subjectSetRepo:         subjectSetRepo,
		}
	}
}

func SubjectSetServiceInstance() *SubjectSetService {
	return subjectSetService
}

func (s *SubjectSetService) Create(ctx context.Context, req *valueobject.CreateSubjectSetReq) (*entity.SubjectSet, *errors.DomainError) {
	// get service
	service, infraErr := s.serviceRepo.GetByNamespace(ctx, req.Namespace)
	if infraErr != nil {
		return nil, errt.DomainTransformerInstance().InfraErrToDomainErr(infraErr)
	}

	// get relation definition
	definition, infraErr := s.relationDefinitionRepo.GetByNamespaceAndRelation(ctx, req.Namespace, req.Relation)
	if infraErr != nil {
		return nil, errt.DomainTransformerInstance().InfraErrToDomainErr(infraErr)
	}

	// create
	set := &entity.SubjectSet{
		ServiceID:            service.ID,
		RelationDefinitionID: definition.ID,
		Namespace:            service.Namespace,
		Object:               req.Object,
		Relation:             definition.Relation,
		Description:          req.Description,
		BaseCreatedUpdated: entity.BaseCreatedUpdated{
			CreatedBy: req.CreatedBy,
			UpdatedBy: req.CreatedBy,
		},
	}
	set, infraErr = s.subjectSetRepo.Create(ctx, set)
	if infraErr != nil {
		return nil, errt.DomainTransformerInstance().InfraErrToDomainErr(infraErr)
	}

	return set, nil
}
