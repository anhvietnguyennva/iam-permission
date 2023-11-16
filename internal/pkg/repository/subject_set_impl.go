package repository

import (
	"context"
	e "errors"
	"strings"

	"github.com/anhvietnguyennva/go-error/pkg/errors"
	"gorm.io/gorm"

	"iam-permission/internal/pkg/constant"
	"iam-permission/internal/pkg/entity"
	"iam-permission/internal/pkg/lock"
	"iam-permission/internal/pkg/repository/postgres"
	"iam-permission/internal/pkg/util/logger"
)

type SubjectSetRepository struct {
	db *gorm.DB
}

var subjectSetRepository *SubjectSetRepository

func InitSubjectSetRepository(db *gorm.DB) {
	lock.InitComponentLock.Lock()
	defer lock.InitComponentLock.Unlock()
	if subjectSetRepository == nil {
		subjectSetRepository = &SubjectSetRepository{
			db: db,
		}
	}
}

func SubjectSetRepositoryInstance() *SubjectSetRepository {
	return subjectSetRepository
}

func (r *SubjectSetRepository) Create(ctx context.Context, set *entity.SubjectSet) (*entity.SubjectSet, *errors.InfraError) {
	model := new(postgres.SubjectSet).FromEntity(set)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		infraErr := errors.NewInfraErrorDBInsert(err, constant.FieldSubjectSet)
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			infraErr = errors.NewInfraErrorDBDuplicatedKey(err, constant.FieldSubjectSet)
		}
		return nil, infraErr
	}

	return model.ToEntity(), nil
}

func (r *SubjectSetRepository) GetByNamespaceAndObjectAndRelation(ctx context.Context, namespace string, object string, relation string) (*entity.SubjectSet, *errors.InfraError) {
	var model postgres.SubjectSet
	if err := r.db.WithContext(ctx).
		Where("namespace = ?", namespace).
		Where("object = ?", object).
		Where("relation = ?", relation).
		First(&model).Error; err != nil {
		infraErr := errors.NewInfraErrorDBSelect(err, constant.FieldSubjectSet)
		if e.Is(err, gorm.ErrRecordNotFound) {
			infraErr = errors.NewInfraErrorDBNotFound(err, constant.FieldSubjectSet)
		}
		logger.Error(ctx, infraErr)
		return nil, infraErr
	}
	return model.ToEntity(), nil
}
