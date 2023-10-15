package postgres

import "iam-permission/internal/pkg/entity"

type Service struct {
	BaseID
	BaseCreatedUpdated

	Namespace   string `gorm:"<-:create"`
	Description string
}

func (s *Service) TableName() string {
	return "services"
}

func (s *Service) FromEntity(e *entity.Service) *Service {
	s.ID = e.ID
	s.Namespace = e.Namespace
	s.Description = e.Description
	s.BaseCreatedUpdated = BaseCreatedUpdated{}.FromEntity(e.BaseCreatedUpdated)
	return s
}

func (s *Service) ToEntity() *entity.Service {
	return &entity.Service{
		BaseID:             s.BaseID.ToEntity(),
		BaseCreatedUpdated: s.BaseCreatedUpdated.ToEntity(),
		Namespace:          s.Namespace,
		Description:        s.Description,
	}
}
