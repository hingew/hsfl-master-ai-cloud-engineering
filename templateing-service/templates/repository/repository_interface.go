package repository

import (
	"github.com/hingew/hsfl-master-ai-cloud-engineering/templating-service/templates/model"
)

type IRepository interface {
	GetAllTemplates() ([]*model.PdfTemplate, error)
	GetTemplate(id int64) (*model.PdfTemplate, error)
	CreateTemplate(template *model.PdfTemplate) error
	UpdateTemplate(id int64, template *model.PdfTemplate) error
	DeleteTemplate(id int64) error
}
