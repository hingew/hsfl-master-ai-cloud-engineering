package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/hingew/hsfl-master-ai-cloud-engineering/templating-service/api/repository"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/templating-service/templates/model"
)

type Controller struct {
	repo repository.IRepository
}

func (c *Controller) GetAllTemplates(w http.ResponseWriter, r *http.Request) {
	templates, err := c.repo.GetAllTemplates()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(templates)
}

func (c *Controller) GetTemplate(w http.ResponseWriter, r *http.Request) {
	templateId := r.Context().Value("id").(string)

	id, err := strconv.ParseInt(templateId, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	template, err := c.repo.GetTemplate(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(template)
}

func (c *Controller) CreateTemplate(w http.ResponseWriter, r *http.Request) {
	var request *model.PdfTemplateCreationRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !request.IsValid() {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := c.repo.CreateTemplate(request); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c *Controller) UpdateTemplate(w http.ResponseWriter, r *http.Request) {
	templateId := r.Context().Value("id").(string)

	id, err := strconv.ParseInt(templateId, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var request *model.PdfTemplateCreationRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := c.repo.UpdateTemplate(id, request); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c *Controller) DeleteTemplate(w http.ResponseWriter, r *http.Request) {
	templateId := r.Context().Value("id").(string)

	id, err := strconv.ParseInt(templateId, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := c.repo.DeleteTemplate(id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
