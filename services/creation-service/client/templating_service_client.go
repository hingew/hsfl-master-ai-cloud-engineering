package client

import "github.com/hingew/hsfl-master-ai-cloud-engineering/lib/model"

type TemplatingServiceClientInterface interface {
	FetchTemplate(templateID uint) (*model.PdfTemplate, error)
}
