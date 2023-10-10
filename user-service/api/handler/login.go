package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

func (r *loginRequest) isValid() bool {
	return r.Email != "" && r.Password != ""
}

type LoginHandler struct {
	userRepository user.Repository
	hasher         crypto.Hasher
	tokenGenerator auth.TokenGenerator
}

func NewLoginHandler(
	userRepository user.Repository,
	hasher crypto.Hasher,
	tokenGenerator auth.TokenGenerator,
) *LoginHandler {
	return &LoginHandler{userRepository, hasher, tokenGenerator}
}

func (handler *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var request loginRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if !request.isValid() {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		users, err := handler.userRepository.FindByEmail(request.Email)
		if err != nil {
			log.Printf("could not find user by email: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if len(users) < 1 {
			w.Header().Add("WWW-Authenticate", "Basic realm=Restricted")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if ok := handler.hasher.Validate([]byte(request.Password), users[0].Password); !ok {
			w.Header().Add("WWW-Authenticate", "Basic realm=Restricted")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		expiration := 1 * time.Hour
		accessToken, err := handler.tokenGenerator.CreateToken(map[string]interface{}{
			"email": request.Email,
			"exp":   time.Now().Add(expiration).Unix(),
		})

		json.NewEncoder(w).Encode(loginResponse{
			AccessToken: accessToken,
			TokenType:   "Bearer",
			ExpiresIn:   int(expiration.Seconds()),
		})
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
