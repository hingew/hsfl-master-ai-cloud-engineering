package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"golang.org/x/sync/singleflight"

	"github.com/hingew/hsfl-master-ai-cloud-engineering/lib/model"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/templateing-service/templates/repository"
)

type ControllerImp struct {
	repo    repository.Repository
	sfGroup *singleflight.Group
}

type createResponse struct {
	Id uint `json:"id"`
}

func NewController(
	repo repository.Repository,
) *ControllerImp {
	g := &singleflight.Group{}
	return &ControllerImp{repo, g}
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

func (c *ControllerImp) GetTemplateWithCoalecing(w http.ResponseWriter, r *http.Request) {
	id, err := c.extractId(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	template, err, _ := c.sfGroup.Do(r.URL.Path, func() (interface{}, error) {
		return c.repo.GetTemplateById(*id)
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.writeTemplateAsResponse(w, template)
}

func (c *ControllerImp) GetTemplate(w http.ResponseWriter, r *http.Request) {
	id, err := c.extractId(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	template, err := c.repo.GetTemplateById(*id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.writeTemplateAsResponse(w, template)
}

func (c *ControllerImp) extractId(r *http.Request) (*uint, error) {
	templateId := r.Context().Value("id").(string)

	id_, err := strconv.Atoi(templateId)
	id := uint(id_)
	if err != nil {
		return nil, err
	}

	return &id, nil
}

func (c *ControllerImp) writeTemplateAsResponse(w http.ResponseWriter, template interface{}) {
	pdfTemplate, ok := template.(*model.PdfTemplate)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pdfTemplate)
}

func (c *ControllerImp) CreateTemplate(w http.ResponseWriter, r *http.Request) {
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

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(createResponse{
		Id: *id,
	})

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
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
