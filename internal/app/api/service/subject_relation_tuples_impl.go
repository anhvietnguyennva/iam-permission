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

type SubjectRelationTupleService struct {
	serviceRepo                    repo.IServiceRepository
	relationDefinitionRepo         repo.IRelationDefinition
	subjectRelationTupleRepository repo.ISubjectRelationTupleRepository
}

var subjectRelationTupleService *SubjectRelationTupleService

func InitSubjectRelationTupleService(
	serviceRepo repo.IServiceRepository,
	relationDefinitionRepo repo.IRelationDefinition,
	subjectRelationTupleRepository repo.ISubjectRelationTupleRepository,
) {
	lock.InitComponentLock.Lock()
	defer lock.InitComponentLock.Unlock()
	if subjectRelationTupleService == nil {
		subjectRelationTupleService = &SubjectRelationTupleService{
			serviceRepo:                    serviceRepo,
			relationDefinitionRepo:         relationDefinitionRepo,
			subjectRelationTupleRepository: subjectRelationTupleRepository,
		}
	}
}

func SubjectRelationTupleServiceInstance() *SubjectRelationTupleService {
	return subjectRelationTupleService
}

func (s *SubjectRelationTupleService) Create(ctx context.Context, req *valueobject.CreateSubjectRelationTupleReq) (*entity.SubjectRelationTuple, *errors.DomainError) {
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
	tuple := &entity.SubjectRelationTuple{
		ServiceID:            service.ID,
		RelationDefinitionID: definition.ID,
		Namespace:            service.Namespace,
		Object:               req.Object,
		Relation:             definition.Relation,
		SubjectID:            req.SubjectID,
		BaseCreatedUpdated: entity.BaseCreatedUpdated{
			CreatedBy: req.CreatedBy,
			UpdatedBy: req.CreatedBy,
		},
	}
	tuple, infraErr = s.subjectRelationTupleRepository.Create(ctx, tuple)
	if infraErr != nil {
		return nil, errt.DomainTransformerInstance().InfraErrToDomainErr(infraErr)
	}

	return tuple, nil
}
