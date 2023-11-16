package repository

import (
	"context"

	"github.com/anhvietnguyennva/go-error/pkg/errors"
	"gorm.io/gorm"

	"iam-permission/internal/pkg/constant"
	"iam-permission/internal/pkg/lock"
	"iam-permission/internal/pkg/repository/postgres"
	"iam-permission/internal/pkg/util/logger"
)

type RelationConfigurationRepository struct {
	db *gorm.DB
}

var relationConfigurationRepository *RelationConfigurationRepository

func InitRelationConfigurationRepository(db *gorm.DB) {
	lock.InitComponentLock.Lock()
	defer lock.InitComponentLock.Unlock()
	if relationConfigurationRepository == nil {
		relationConfigurationRepository = &RelationConfigurationRepository{
			db: db,
		}
	}
}

func RelationConfigurationRepositoryInstance() *RelationConfigurationRepository {
	return relationConfigurationRepository
}

func (r *RelationConfigurationRepository) GetChildrenRelations(ctx context.Context, namespace string, parentRelation string) ([]string, *errors.InfraError) {
	var childrenRelations []string
	if err := r.db.WithContext(ctx).
		Model(&postgres.RelationConfiguration{}).
		Where("namespace = ?", namespace).
		Where("parent_relation = ?", parentRelation).
		Select("child_relation").
		Scan(&childrenRelations).
		Error; err != nil {
		infraErr := errors.NewInfraErrorDBSelect(err, constant.FieldRelationConfiguration)
		logger.Error(ctx, infraErr)
		return nil, infraErr
	}
	return childrenRelations, nil
}
