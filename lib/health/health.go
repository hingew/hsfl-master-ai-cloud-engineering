package health

import (
	"encoding/json"
	"net/http"
)

type Health struct {
	status string
}

func Check(w http.ResponseWriter, r *http.Request) {
	health := Health{status: "ok"}
	json.NewEncoder(w).Encode(health)
	w.WriteHeader(http.StatusOK)
}
