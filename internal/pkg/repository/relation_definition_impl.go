package repository

import (
	"context"
	e "errors"

	"github.com/anhvietnguyennva/go-error/pkg/errors"
	"gorm.io/gorm"

	"iam-permission/internal/pkg/constant"
	"iam-permission/internal/pkg/entity"
	"iam-permission/internal/pkg/lock"
	"iam-permission/internal/pkg/repository/postgres"
	"iam-permission/internal/pkg/util/logger"
)

type RelationDefinitionRepository struct {
	db *gorm.DB
}

var relationDefinitionRepository *RelationDefinitionRepository

func InitRelationDefinitionRepository(db *gorm.DB) {
	lock.InitComponentLock.Lock()
	defer lock.InitComponentLock.Unlock()
	if relationDefinitionRepository == nil {
		relationDefinitionRepository = &RelationDefinitionRepository{
			db: db,
		}
	}
}

func RelationDefinitionRepositoryInstance() *RelationDefinitionRepository {
	return relationDefinitionRepository
}

func (r *RelationDefinitionRepository) GetByNamespaceAndRelation(ctx context.Context, namespace string, relation string) (*entity.RelationDefinition, *errors.InfraError) {
	var model postgres.RelationDefinition
	if err := r.db.WithContext(ctx).
		Where("namespace = ?", namespace).
		Where("relation = ?", relation).
		First(&model).Error; err != nil {
		infraErr := errors.NewInfraErrorDBSelect(err, constant.FieldRelationDefinition)
		if e.Is(err, gorm.ErrRecordNotFound) {
			infraErr = errors.NewInfraErrorDBNotFound(err, constant.FieldRelationDefinition)
		}
		logger.Error(ctx, infraErr)
		return nil, infraErr
	}
	return model.ToEntity(), nil
}
