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
	serviceModel := new(postgres.Service).FromEntity(service)
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// create service
		if err := tx.Create(&serviceModel).Error; err != nil {
			infraErr := errors.NewInfraErrorDBInsert(err, constant.FieldService)
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				infraErr = errors.NewInfraErrorDBDuplicatedKey(err, constant.FieldService)
			}
			return infraErr
		}

		// create default relations
		relationViewer := &postgres.RelationDefinition{
			ServiceID:   serviceModel.ID,
			Namespace:   serviceModel.Namespace,
			Relation:    constant.RelationViewer,
			Description: constant.RelationViewer,
			BaseCreatedUpdated: postgres.BaseCreatedUpdated{
				CreatedBy: serviceModel.CreatedBy,
				UpdatedBy: serviceModel.CreatedBy,
			},
		}
		relationEditor := &postgres.RelationDefinition{
			ServiceID:   serviceModel.ID,
			Namespace:   serviceModel.Namespace,
			Relation:    constant.RelationEditor,
			Description: constant.RelationEditor,
			BaseCreatedUpdated: postgres.BaseCreatedUpdated{
				CreatedBy: serviceModel.CreatedBy,
				UpdatedBy: serviceModel.CreatedBy,
			},
		}
		relationOwner := &postgres.RelationDefinition{
			ServiceID:   serviceModel.ID,
			Namespace:   serviceModel.Namespace,
			Relation:    constant.RelationOwner,
			Description: constant.RelationOwner,
			BaseCreatedUpdated: postgres.BaseCreatedUpdated{
				CreatedBy: serviceModel.CreatedBy,
				UpdatedBy: serviceModel.CreatedBy,
			},
		}
		relationConsumer := &postgres.RelationDefinition{
			ServiceID:   serviceModel.ID,
			Namespace:   serviceModel.Namespace,
			Relation:    constant.RelationConsumer,
			Description: constant.RelationConsumer,
			BaseCreatedUpdated: postgres.BaseCreatedUpdated{
				CreatedBy: serviceModel.CreatedBy,
				UpdatedBy: serviceModel.CreatedBy,
			},
		}
		defaultRelationDefinitions := []*postgres.RelationDefinition{relationViewer, relationEditor, relationOwner, relationConsumer}
		if err := tx.CreateInBatches(defaultRelationDefinitions, constant.DefaultDBBatchSize).Error; err != nil {
			return errors.NewInfraErrorDBInsert(err, constant.FieldRelationDefinition)
		}

		// create default relation configurations
		configurationViewerEditor := &postgres.RelationConfiguration{
			ServiceID:                  serviceModel.ID,
			Namespace:                  serviceModel.Namespace,
			ParentRelationDefinitionID: relationViewer.ID,
			ParentRelation:             relationViewer.Relation,
			ChildRelationDefinitionID:  relationEditor.ID,
			ChildRelation:              relationEditor.Relation,
			BaseCreatedUpdated: postgres.BaseCreatedUpdated{
				CreatedBy: serviceModel.CreatedBy,
				UpdatedBy: serviceModel.CreatedBy,
			},
		}
		configurationEditorOwner := &postgres.RelationConfiguration{
			ServiceID:                  serviceModel.ID,
			Namespace:                  serviceModel.Namespace,
			ParentRelationDefinitionID: relationEditor.ID,
			ParentRelation:             relationEditor.Relation,
			ChildRelationDefinitionID:  relationOwner.ID,
			ChildRelation:              relationOwner.Relation,
			BaseCreatedUpdated: postgres.BaseCreatedUpdated{
				CreatedBy: serviceModel.CreatedBy,
				UpdatedBy: serviceModel.CreatedBy,
			},
		}
		defaultRelationConfigurations := []*postgres.RelationConfiguration{configurationViewerEditor, configurationEditorOwner}
		if err := tx.CreateInBatches(defaultRelationConfigurations, constant.DefaultDBBatchSize).Error; err != nil {
			return errors.NewInfraErrorDBInsert(err, constant.FieldRelationDefinition)
		}

		return nil
	})

	if err != nil {
		var infraErr *errors.InfraError
		e.As(err, &infraErr)
		logger.Error(ctx, infraErr)
		return nil, infraErr
	}
	return serviceModel.ToEntity(), nil
}
