package controller

import "net/http"

type ControllerInterface interface {
	Login(http.ResponseWriter, *http.Request)
	Register(http.ResponseWriter, *http.Request)
}
