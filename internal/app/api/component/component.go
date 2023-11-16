package component

import (
	"iam-permission/internal/app/api/service"
	"iam-permission/internal/pkg/db"
	"iam-permission/internal/pkg/redis"
	repo "iam-permission/internal/pkg/repository"
	"iam-permission/internal/pkg/util/logger"
)

func InitComponents() error {
	// Logger
	if err := logger.InitLogger(); err != nil {
		return err
	}

	// Errors
	registerAPIServerErrors()

	// Validators
	if err := registerValidators(); err != nil {
		return err
	}

	// Infra
	if err := db.InitDB(); err != nil {
		return err
	}
	if err := redis.InitClient(); err != nil {
		return err
	}

	// Repo
	repo.InitServiceRepository(db.Instance())
	repo.InitRelationDefinitionRepository(db.Instance())
	repo.InitRelationDefinitionRepository(db.Instance())
	repo.InitSubjectRelationTupleRepository(db.Instance())
	repo.InitSubjectSetRepository(db.Instance())
	repo.InitSubjectSetRelationTupleRepository(db.Instance())
	repo.InitRelationConfigurationRepository(db.Instance())

	// Service
	service.InitServiceService(repo.ServiceRepositoryInstance())
	service.InitRelationDefinitionService(repo.ServiceRepositoryInstance(), repo.RelationDefinitionRepositoryInstance())
	service.InitSubjectRelationTupleService(
		repo.ServiceRepositoryInstance(),
		repo.RelationDefinitionRepositoryInstance(),
		repo.SubjectRelationTupleRepositoryInstance(),
	)
	service.InitSubjectSetService(
		repo.ServiceRepositoryInstance(),
		repo.RelationDefinitionRepositoryInstance(),
		repo.SubjectSetRepositoryInstance(),
	)
	service.InitSubjectSetRelationTupleService(
		repo.ServiceRepositoryInstance(),
		repo.RelationDefinitionRepositoryInstance(),
		repo.SubjectSetRepositoryInstance(),
		repo.SubjectSetRelationTupleRepositoryInstance(),
	)
	service.InitPermissionService(
		repo.SubjectRelationTupleRepositoryInstance(),
		repo.SubjectSetRepositoryInstance(),
		repo.SubjectSetRelationTupleRepositoryInstance(),
		repo.RelationConfigurationRepositoryInstance(),
	)

	return nil
}

func registerAPIServerErrors() {
	// Validation errors

	// Infra errors
}

func registerValidators() error {
	//v, ok := binding.Validator.Engine().(*validator.Validate)
	//if ok {
	//}

	return nil
}
