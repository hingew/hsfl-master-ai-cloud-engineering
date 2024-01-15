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
	templatingClient client.TemplatingServiceClient
}

func NewController(pdf pdf.Pdf, templatingClient client.TemplatingServiceClient) *Controller {
	return &Controller{pdf, templatingClient}
}

func isValid(template *model.PdfTemplate, params map[string]string) bool {
	if len(template.Elements) == 0 {
		return true
	}

	for _, element := range template.Elements {
		if element.Type == "text" && element.ValueFrom != "" {
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

	template, err := c.templatingClient.FetchTemplate(id)

	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var request map[string]string
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !isValid(template, request) {
		log.Print("IS not valid: ", request)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	report := pdf.New()
	report.Render(template, request)
	buf, err := report.Out()

	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-type", "application/pdf")
	w.Write(buf.Bytes())

}
