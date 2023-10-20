package controller

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/hingew/hsfl-master-ai-cloud-engineering/creation-service/pdf"
)

type Controller struct {
	pdf pdf.Pdf
}

func NewController(pdf pdf.Pdf) *Controller {
	return &Controller{pdf}
}

type creationRequest struct {
	TemplateId uint
}

func (r *creationRequest) isValid() bool {
	//TODO: add request validation
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

	var request creationRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !request.isValid() {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	report := pdf.New()

	// TODO: render elements

	buf, err := report.Out()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/pdf")

	//Stream to response
	if _, err := io.Copy(w, buf); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

}
