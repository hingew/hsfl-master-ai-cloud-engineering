package repository

import "github.com/hingew/hsfl-master-ai-cloud-engineering/user-service/user/model"

type RepositoryInterface interface {
	Migrate() error
	Create([]*model.DbUser) error
	FindByEmail(email string) ([]*model.DbUser, error)
	Delete([]*model.DbUser) error
}
