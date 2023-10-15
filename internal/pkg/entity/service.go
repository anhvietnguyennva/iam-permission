package entity

type Service struct {
	BaseID
	BaseCreatedUpdated
	Namespace   string
	Description string
}
