package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"golang.org/x/sync/singleflight"

	"github.com/hingew/hsfl-master-ai-cloud-engineering/lib/model"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/templateing-service/templates/repository"
)

type Controller struct {
	repo    repository.Repository
	sfGroup *singleflight.Group
}

type createResponse struct {
	Id uint `json:"id"`
}

func NewController(
	repo repository.Repository,
) *Controller {
	g := &singleflight.Group{}
	return &Controller{repo, g}
}

func (c *Controller) GetAllTemplates(w http.ResponseWriter, r *http.Request) {
	templates, err := c.repo.GetAllTemplates()
	if err != nil {
        log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(templates)
}

func (c *Controller) GetTemplateWithCoalecing(w http.ResponseWriter, r *http.Request) {
	id, err := c.extractId(r)
	if err != nil {
        log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	template, err, _ := c.sfGroup.Do(r.URL.Path, func() (interface{}, error) {
		return c.repo.GetTemplateById(*id)
	})
	if err != nil {
        log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.writeTemplateAsResponse(w, template)
}

func (c *Controller) GetTemplate(w http.ResponseWriter, r *http.Request) {
	id, err := c.extractId(r)
	if err != nil {
        log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	template, err := c.repo.GetTemplateById(*id)
	if err != nil {
        log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.writeTemplateAsResponse(w, template)
}

func (c *Controller) extractId(r *http.Request) (*uint, error) {
	templateId := r.Context().Value("id").(string)

	id_, err := strconv.Atoi(templateId)
	id := uint(id_)
	if err != nil {
        log.Println(err)
		return nil, err
	}

	return &id, nil
}

func (c *Controller) writeTemplateAsResponse(w http.ResponseWriter, template interface{}) {
	pdfTemplate, ok := template.(*model.PdfTemplate)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pdfTemplate)
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

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(createResponse{
		Id: *id,
	})

}

func (c *Controller) UpdateTemplate(w http.ResponseWriter, r *http.Request) {
	templateId := r.Context().Value("id").(string)

	id_, err := strconv.Atoi(templateId)
	id := uint(id_)
	if err != nil {
        log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var request *model.PdfTemplate
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
        log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := c.repo.UpdateTemplate(id, *request); err != nil {
        log.Println(err)
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
        log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := c.repo.DeleteTemplate(id); err != nil {
        log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
