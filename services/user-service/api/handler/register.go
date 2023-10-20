package handler

import (
	"encoding/json"
	"net/http"

	"github.com/hingew/hsfl-master-ai-cloud-engineering/user-service/crypto"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/user-service/user"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/user-service/user/model"
)

type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *registerRequest) isValid() bool {
	return r.Email != "" && r.Password != ""
}

func Register(
	userRepository user.Repository,
	hasher crypto.Hasher,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request registerRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if !request.isValid() {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		products, err := userRepository.FindByEmail(request.Email)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if len(products) > 0 {
			w.WriteHeader(http.StatusConflict)
			return
		}

		hashedPassword, err := hasher.Hash([]byte(request.Password))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err := userRepository.Create([]*model.DbUser{{
			Email:    request.Email,
			Password: hashedPassword,
		}}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
