package user

import "github.com/hingew/hsfl-master-ai-cloud-engineering/user-serivce/user/model"

type Repository interface {
	Migrate() error
	Create([]*model.DbUser) error
	FindByEmail(email string) ([]*model.DbUser, error)
	Delete([]*model.DbUser) error
}
