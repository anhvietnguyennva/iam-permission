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

type SubjectSetRelationTupleRepository struct {
	db *gorm.DB
}

var subjectSetRelationTupleRepository *SubjectSetRelationTupleRepository

func InitSubjectSetRelationTupleRepository(db *gorm.DB) {
	lock.InitComponentLock.Lock()
	defer lock.InitComponentLock.Unlock()
	if subjectSetRelationTupleRepository == nil {
		subjectSetRelationTupleRepository = &SubjectSetRelationTupleRepository{
			db: db,
		}
	}
}

func SubjectSetRelationTupleRepositoryInstance() *SubjectSetRelationTupleRepository {
	return subjectSetRelationTupleRepository
}

func (r *SubjectSetRelationTupleRepository) Create(ctx context.Context, tuple *entity.SubjectSetRelationTuple) (*entity.SubjectSetRelationTuple, *errors.InfraError) {
	model := new(postgres.SubjectSetRelationTuple).FromEntity(tuple)

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		infraErr := errors.NewInfraErrorDBInsert(err, constant.FieldSubjectSetRelationTuple)
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			infraErr = errors.NewInfraErrorDBDuplicatedKey(err, constant.FieldSubjectSetRelationTuple)
		}
		logger.Error(ctx, infraErr)
		return nil, infraErr
	}

	return model.ToEntity(), nil
}
