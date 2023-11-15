package repository

import (
	"context"
	"strings"

	"github.com/anhvietnguyennva/go-error/pkg/errors"
	"gorm.io/gorm"

	"iam-permission/internal/pkg/constant"
	"iam-permission/internal/pkg/entity"
	"iam-permission/internal/pkg/lock"
	"iam-permission/internal/pkg/repository/postgres"
	"iam-permission/internal/pkg/util/logger"
)

type SubjectRelationTupleRepository struct {
	db *gorm.DB
}

var subjectRelationTupleRepository *SubjectRelationTupleRepository

func InitSubjectRelationTupleRepository(db *gorm.DB) {
	lock.InitComponentLock.Lock()
	defer lock.InitComponentLock.Unlock()
	if subjectRelationTupleRepository == nil {
		subjectRelationTupleRepository = &SubjectRelationTupleRepository{
			db: db,
		}
	}
}

func SubjectRelationTupleRepositoryInstance() *SubjectRelationTupleRepository {
	return subjectRelationTupleRepository
}

func (r *SubjectRelationTupleRepository) Create(ctx context.Context, tuple *entity.SubjectRelationTuple) (*entity.SubjectRelationTuple, *errors.InfraError) {
	model := new(postgres.SubjectRelationTuple).FromEntity(tuple)

	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		infraErr := errors.NewInfraErrorDBInsert(err, constant.FieldSubjectRelationTuple)
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			infraErr = errors.NewInfraErrorDBDuplicatedKey(err, constant.FieldSubjectRelationTuple)
		}
		logger.Error(ctx, infraErr)
		return nil, infraErr
	}

	return model.ToEntity(), nil
}
