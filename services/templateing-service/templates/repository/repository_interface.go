package repository

import (
	"github.com/hingew/hsfl-master-ai-cloud-engineering/templating-service/templates/model"
)

type Repository interface {
	CreateTemplate(data model.PdfTemplate) (*uint, error)
	GetAllTemplates() ([]*model.PdfTemplate, error)
	GetTemplateById(id uint) (*model.PdfTemplate, error)
	UpdateTemplate(id uint, data model.PdfTemplate) error
	DeleteTemplate(id uint) error
}
