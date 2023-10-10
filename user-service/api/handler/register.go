package handler

import (
	"encoding/json"
	"net/http"
)

type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *registerRequest) isValid() bool {
	return r.Email != "" && r.Password != ""
}

type RegisterHandler struct {
	userRepository user.Repository
	hasher         crypto.Hasher
}

func NewRegisterHandler(
	userRepository user.Repository,
	hasher crypto.Hasher,
) *RegisterHandler {
	return &RegisterHandler{userRepository, hasher}
}

func (handler *RegisterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var request registerRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if !request.isValid() {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		products, err := handler.userRepository.FindByEmail(request.Email)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if len(products) > 0 {
			w.WriteHeader(http.StatusConflict)
			return
		}

		hashedPassword, err := handler.hasher.Hash([]byte(request.Password))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err := handler.userRepository.Create([]*model.DbUser{{
			Email:    request.Email,
			Password: hashedPassword,
		}}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
