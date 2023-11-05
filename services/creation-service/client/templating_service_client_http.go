package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/hingew/hsfl-master-ai-cloud-engineering/lib/model"
)

type TemplatingClient struct {
	url string
}

func NewClient(url string) TemplatingClient {
	return TemplatingClient{url}
}

func (t *TemplatingClient) GetTemplate(templateID uint) (*model.PdfTemplate, error) {
	requestURL := fmt.Sprintf("%s/api/templates/%d", t.url, templateID)

	res, err := http.Get(requestURL)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("Templating Service returned %d HTTP status", res.StatusCode))
	}

	var template *model.PdfTemplate
	err = json.NewDecoder(res.Body).Decode(&template)

	return template, err
}
