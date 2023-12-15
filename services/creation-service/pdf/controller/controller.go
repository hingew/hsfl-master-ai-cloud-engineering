package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/hingew/hsfl-master-ai-cloud-engineering/creation-service/client"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/creation-service/pdf"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/lib/model"
)

type Controller struct {
	pdf              pdf.Pdf
	templatingClient client.TemplatingGrpcClient
}

func NewController(pdf pdf.Pdf, templatingClient client.TemplatingGrpcClient) *Controller {
	return &Controller{pdf, templatingClient}
}

func isValid(template *model.PdfTemplate, params map[string]interface{}) bool {
	for _, element := range template.Elements {
		if element.ValueFrom != "" {
			if _, ok := params[element.ValueFrom]; !ok {
				return false
			}
		}
	}
	return true
}

func (c *Controller) CreatePdf(w http.ResponseWriter, r *http.Request) {
	templateId := r.Context().Value("id").(string)

	id_, err := strconv.Atoi(templateId)
	id := uint(id_)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// GET the template from the templating service
	template, err := c.templatingClient.FetchTemplate(id)
	log.Print(template)
	log.Print(err)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !isValid(template, request) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	report := pdf.New()
	report.Render(template, request)
	buf, err := report.Out()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-type", "application/pdf")
	w.Write(buf.Bytes())

}
