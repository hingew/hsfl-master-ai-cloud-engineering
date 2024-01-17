package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/hingew/hsfl-master-ai-cloud-engineering/user-service/auth"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/user-service/crypto"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/user-service/user/model"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/user-service/user/repository"
)

type Controller struct {
	repo           repository.RepositoryInterface
	hasher         crypto.Hasher
	tokenGenerator auth.TokenGenerator
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

type registerRequest struct {
	Email                string `json:"email"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"password_confirmation"`
}

func (r *loginRequest) isValid() bool {
	return r.Email != "" && r.Password != ""
}

func (r *registerRequest) isValid() bool {
	return r.Email != "" && r.Password != "" && r.Password == r.PasswordConfirmation
}

func NewController(repo repository.RepositoryInterface, hasher crypto.Hasher, tokenGenerator auth.TokenGenerator) *Controller {
	return &Controller{repo, hasher, tokenGenerator}
}

func (c *Controller) Login(w http.ResponseWriter, r *http.Request) {
	var request loginRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Printf("Could not decode request: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !request.isValid() {
        log.Println("Request is not valid: ", request)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	users, err := c.repo.FindByEmail(request.Email)
	if err != nil {
		log.Printf("could not find user by email: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(users) < 1 {
		log.Println("User does not exists")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if ok := c.hasher.Validate([]byte(request.Password), users[0].Password); !ok {
        log.Printf("password is incorrect")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	expiration := 1 * time.Hour
	accessToken, err := c.tokenGenerator.CreateToken(map[string]interface{}{
		"email": request.Email,
		"exp":   time.Now().Add(expiration).Unix(),
	})

    err = json.NewEncoder(w).Encode(loginResponse{
		AccessToken: accessToken,
		TokenType:   "Bearer",
		ExpiresIn:   int(expiration.Seconds()),
	})

    if err != nil {
        log.Printf("Encode user failed: %s", err.Error())
    }
}


func (c *Controller) Register(w http.ResponseWriter, r *http.Request) {
	var request registerRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Printf("Failed decoding request: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !request.isValid() {
		log.Println("Request is not valid: ", request)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := c.repo.FindByEmail(request.Email)
	if err != nil {
		log.Printf("Could not find user by email: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(user) > 0 {
		log.Println("User already exists")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	hashedPassword, err := c.hasher.Hash([]byte(request.Password))
	if err != nil {
		log.Printf("Hasing the password faild: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := c.repo.Create([]*model.DbUser{{
		Email:    request.Email,
		Password: hashedPassword,
	}}); err != nil {
		log.Printf("Could not create user: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
