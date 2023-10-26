package controller

import "net/http"

type Controller interface {
	GetAllTemplates(http.ResponseWriter, *http.Request)
	GetTemplate(http.ResponseWriter, *http.Request)
	CreateTemplate(http.ResponseWriter, *http.Request)
	UpdateTemplate(http.ResponseWriter, *http.Request)
	DeleteTemplate(http.ResponseWriter, *http.Request)
}
