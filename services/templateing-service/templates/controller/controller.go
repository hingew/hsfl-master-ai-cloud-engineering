package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/hingew/hsfl-master-ai-cloud-engineering/templating-service/templates/model"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/templating-service/templates/repository"
)

type Controller struct {
	repo repository.IRepository
}

func NewController(
	repo repository.IRepository,
) *Controller {
	return &Controller{repo}
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

	id_, err := strconv.Atoi(templateId)
	id := uint(id_)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	template, err := c.repo.GetTemplateById(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(template)
}

func (c *Controller) CreateTemplate(w http.ResponseWriter, r *http.Request) {
	var request model.PdfTemplate
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := c.repo.CreateTemplate(request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := fmt.Sprintf("id: %d", id)
	w.Write([]byte(response))
	w.WriteHeader(http.StatusOK)
}

func (c *Controller) UpdateTemplate(w http.ResponseWriter, r *http.Request) {
	templateId := r.Context().Value("id").(string)

	id_, err := strconv.Atoi(templateId)
	id := uint(id_)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var request *model.PdfTemplate
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := c.repo.UpdateTemplate(id, *request); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c *Controller) DeleteTemplate(w http.ResponseWriter, r *http.Request) {
	templateId := r.Context().Value("id").(string)

	id_, err := strconv.Atoi(templateId)
	id := uint(id_)
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
