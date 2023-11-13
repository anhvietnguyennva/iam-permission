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

type ServiceRepository struct {
	db *gorm.DB
}

var serviceRepo *ServiceRepository

func InitServiceRepository(db *gorm.DB) {
	lock.InitComponentLock.Lock()
	defer lock.InitComponentLock.Unlock()
	if serviceRepo == nil {
		serviceRepo = &ServiceRepository{
			db: db,
		}
	}
}

func ServiceRepositoryInstance() *ServiceRepository {
	return serviceRepo
}

func (r *ServiceRepository) Create(ctx context.Context, service *entity.Service) (*entity.Service, *errors.InfraError) {
	model := new(postgres.Service).FromEntity(service)
	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		infraErr := errors.NewInfraErrorDBInsert(err, constant.FieldService)
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			infraErr = errors.NewInfraErrorDBDuplicatedKey(err, constant.FieldService)
		}
		logger.Error(ctx, infraErr)
		return nil, infraErr
	}
	return model.ToEntity(), nil
}
