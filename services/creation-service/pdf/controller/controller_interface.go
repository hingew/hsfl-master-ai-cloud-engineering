package controller

import "net/http"

type ControllerInterface interface {
	CreatePdf(http.ResponseWriter, *http.Request)
}

