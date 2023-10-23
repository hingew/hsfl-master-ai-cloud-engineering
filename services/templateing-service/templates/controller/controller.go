package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/hingew/hsfl-master-ai-cloud-engineering/lib/model"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/templateing-service/templates/repository"
)

type ControllerImp struct {
	repo repository.Repository
}

func NewController(
	repo repository.Repository,
) *ControllerImp {
	return &ControllerImp{repo}
}

func (c *ControllerImp) GetAllTemplates(w http.ResponseWriter, r *http.Request) {
	templates, err := c.repo.GetAllTemplates()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(templates)
}

func (c *ControllerImp) GetTemplate(w http.ResponseWriter, r *http.Request) {
	templateId := r.Context().Value("id").(string)

	id_, err := strconv.Atoi(templateId)
	id := uint(id_)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	template, err := c.repo.GetTemplateById(id)
	if err != nil {
		// TODO wenn element nicht vorhanden, http.StatusNotFound
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(template)
}

func (c *ControllerImp) CreateTemplate(w http.ResponseWriter, r *http.Request) {
	var request model.PdfTemplate
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := c.repo.CreateTemplate(request)
	if err != nil {
		// TODO könnte man sich mit etwas mehr aufwand ran setzen --> https://gorm.io/docs/error_handling.html
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := fmt.Sprintf("id: %d", id)
	w.Write([]byte(response))
	w.WriteHeader(http.StatusOK)
}

func (c *ControllerImp) UpdateTemplate(w http.ResponseWriter, r *http.Request) {
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
		// TODO könnte man sich mit etwas mehr aufwand ran setzen --> https://gorm.io/docs/error_handling.html
		// https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/422
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c *ControllerImp) DeleteTemplate(w http.ResponseWriter, r *http.Request) {
	templateId := r.Context().Value("id").(string)

	id_, err := strconv.Atoi(templateId)
	id := uint(id_)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := c.repo.DeleteTemplate(id); err != nil {
		// TODO könnte man sich mit etwas mehr aufwand ran setzen --> https://gorm.io/docs/error_handling.html
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
