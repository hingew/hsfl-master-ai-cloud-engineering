package handler

import (
	"encoding/json"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/creation-service/pdf"
	"io"
	"net/http"
)

type creationRequest struct {
}

func (r *creationRequest) isValid() bool {
	//TODO: add request validation
	return true
}

func Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
}
