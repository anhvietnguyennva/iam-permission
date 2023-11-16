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

type SubjectSetRelationTupleService struct {
	serviceRepo                 repo.IServiceRepository
	relationDefinitionRepo      repo.IRelationDefinitionRepository
	subjectSetRepo              repo.ISubjectSetRepository
	subjectSetRelationTupleRepo repo.ISubjectSetRelationTupleRepository
}

var subjectSetRelationTupleService *SubjectSetRelationTupleService

func InitSubjectSetRelationTupleService(
	serviceRepo repo.IServiceRepository,
	relationDefinitionRepo repo.IRelationDefinitionRepository,
	subjectSetRepo repo.ISubjectSetRepository,
	subjectSetRelationTupleRepository repo.ISubjectSetRelationTupleRepository,
) {
	lock.InitComponentLock.Lock()
	defer lock.InitComponentLock.Unlock()
	if subjectSetRelationTupleService == nil {
		subjectSetRelationTupleService = &SubjectSetRelationTupleService{
			serviceRepo:                 serviceRepo,
			relationDefinitionRepo:      relationDefinitionRepo,
			subjectSetRepo:              subjectSetRepo,
			subjectSetRelationTupleRepo: subjectSetRelationTupleRepository,
		}
	}
}

func SubjectSetRelationTupleServiceInstance() *SubjectSetRelationTupleService {
	return subjectSetRelationTupleService
}

func (s *SubjectSetRelationTupleService) Create(ctx context.Context, req *valueobject.CreateSubjectSetRelationTupleReq) (*entity.SubjectSetRelationTuple, *errors.DomainError) {
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

	// get subject set
	subjectSet, infraErr := s.subjectSetRepo.GetByNamespaceAndObjectAndRelation(ctx, req.SubjectSetNamespace, req.SubjectSetObject, req.SubjectSetRelation)
	if infraErr != nil {
		return nil, errt.DomainTransformerInstance().InfraErrToDomainErr(infraErr)
	}

	// create
	tuple := &entity.SubjectSetRelationTuple{
		ServiceID:            service.ID,
		RelationDefinitionID: definition.ID,
		SubjectSetID:         subjectSet.ID,
		Namespace:            service.Namespace,
		Object:               req.Object,
		Relation:             definition.Relation,
		SubjectSetNamespace:  subjectSet.Namespace,
		SubjectSetObject:     subjectSet.Object,
		SubjectSetRelation:   subjectSet.Relation,
		BaseCreatedUpdated: entity.BaseCreatedUpdated{
			CreatedBy: req.CreatedBy,
			UpdatedBy: req.CreatedBy,
		},
	}
	tuple, infraErr = s.subjectSetRelationTupleRepo.Create(ctx, tuple)
	if infraErr != nil {
		return nil, errt.DomainTransformerInstance().InfraErrToDomainErr(infraErr)
	}

	return tuple, nil
}
