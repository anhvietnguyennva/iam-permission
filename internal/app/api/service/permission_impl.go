package service

import (
	"context"

	errc "github.com/anhvietnguyennva/go-error/pkg/constant"
	"github.com/anhvietnguyennva/go-error/pkg/errors"
	errt "github.com/anhvietnguyennva/go-error/pkg/transformer"

	"iam-permission/internal/app/api/valueobject"
	"iam-permission/internal/pkg/lock"
	repo "iam-permission/internal/pkg/repository"
)

type PermissionService struct {
	subjectRelationTupleRepo    repo.ISubjectRelationTupleRepository
	subjectSetRepo              repo.ISubjectSetRepository
	subjectSetRelationTupleRepo repo.ISubjectSetRelationTupleRepository
	relationConfigurationRepo   repo.IRelationConfigurationRepository
}

var permissionService *PermissionService

func InitPermissionService(
	subjectRelationTupleRepo repo.ISubjectRelationTupleRepository,
	subjectSetRepo repo.ISubjectSetRepository,
	subjectSetRelationTupleRepo repo.ISubjectSetRelationTupleRepository,
	relationConfigurationRepo repo.IRelationConfigurationRepository,
) {
	lock.InitComponentLock.Lock()
	defer lock.InitComponentLock.Unlock()
	if permissionService == nil {
		permissionService = &PermissionService{
			subjectRelationTupleRepo:    subjectRelationTupleRepo,
			subjectSetRepo:              subjectSetRepo,
			subjectSetRelationTupleRepo: subjectSetRelationTupleRepo,
			relationConfigurationRepo:   relationConfigurationRepo,
		}
	}
}

func PermissionServiceInstance() *PermissionService {
	return permissionService
}

func (s *PermissionService) CheckPermission(ctx context.Context, req *valueobject.CheckPermissionReq) (bool, *errors.DomainError) {
	return s.checkPermission(ctx, req, 0)
}

func (s *PermissionService) checkPermission(ctx context.Context, req *valueobject.CheckPermissionReq, depth uint8) (bool, *errors.DomainError) {
	if depth >= req.MaxDepth {
		return false, nil
	}

	// check subject relation tuple
	_, infraErr := s.subjectRelationTupleRepo.GetByNamespaceAndObjectAndRelationAndSubjectID(ctx, req.Namespace, req.Object, req.Relation, req.SubjectID)
	if infraErr != nil && infraErr.Code != errc.InfraErrCodeDBNotFound {
		return false, errt.DomainTransformerInstance().InfraErrToDomainErr(infraErr)
	}
	if infraErr == nil {
		return true, nil
	}

	// check subject set
	subjectSetRelationTuples, infraErr := s.subjectSetRelationTupleRepo.GetByNamespaceAndObjectAndRelation(ctx, req.Namespace, req.Object, req.Relation)
	if infraErr != nil {
		return false, errt.DomainTransformerInstance().InfraErrToDomainErr(infraErr)
	}
	for _, tuple := range subjectSetRelationTuples {
		r := &valueobject.CheckPermissionReq{
			Namespace: tuple.SubjectSetNamespace,
			Object:    tuple.SubjectSetObject,
			Relation:  tuple.SubjectSetRelation,
			SubjectID: req.SubjectID,
			MaxDepth:  req.MaxDepth,
		}
		allowed, domainErr := s.checkPermission(ctx, r, depth+1)
		if domainErr != nil {
			return false, domainErr
		}
		if allowed {
			return allowed, nil
		}
	}

	// check child relation
	parentRelations, infraErr := s.relationConfigurationRepo.GetChildrenRelations(ctx, req.Namespace, req.Relation)
	if infraErr != nil {
		return false, errt.DomainTransformerInstance().InfraErrToDomainErr(infraErr)
	}
	for _, parentRelation := range parentRelations {
		r := &valueobject.CheckPermissionReq{
			Namespace: req.Namespace,
			Object:    req.Object,
			Relation:  parentRelation,
			SubjectID: req.SubjectID,
			MaxDepth:  req.MaxDepth,
		}
		allowed, domainErr := s.checkPermission(ctx, r, depth+1)
		if domainErr != nil {
			return false, domainErr
		}
		if allowed {
			return allowed, nil
		}
	}

	return false, nil
}
