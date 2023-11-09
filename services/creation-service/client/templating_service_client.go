package client

import "github.com/hingew/hsfl-master-ai-cloud-engineering/lib/model"

type TemplatingServiceClient interface {
	GetTemplate(templateID uint) (*model.PdfTemplate, error)
}