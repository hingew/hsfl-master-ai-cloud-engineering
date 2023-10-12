package pdftemplates

import "net/http"

type Controller interface {
	GetTemplates(http.ResponseWriter, *http.Request)
	PostTemplates(http.ResponseWriter, *http.Request)
	GetTemplate(http.ResponseWriter, *http.Request)
	PutTemplates(http.ResponseWriter, *http.Request)
	DeleteTemplates(http.ResponseWriter, *http.Request)
}
