package controller

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/hingew/hsfl-master-ai-cloud-engineering/lib/gomockhelpers"
	mocks "github.com/hingew/hsfl-master-ai-cloud-engineering/user-service/_mocks"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/user-service/user/model"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestValidation(t *testing.T) {
	t.Run("Login Validation", func(t *testing.T) {

		t.Run("Should fail if the email is empty", func(t *testing.T) {
			r := loginRequest{Email: "", Password: "test"}

			assert.False(t, r.isValid())
		})

		t.Run("Should fail if the password is empty", func(t *testing.T) {
			r := loginRequest{Email: "test@test.de", Password: ""}

			assert.False(t, r.isValid())
		})

		t.Run("Should succeed if fields are correct", func(t *testing.T) {
			r := loginRequest{Email: "test@test.de", Password: "test"}

			assert.True(t, r.isValid())
		})
	})

	t.Run("Register Validation", func(t *testing.T) {

		t.Run("Should fail if the email is empty", func(t *testing.T) {
			r := registerRequest{Email: "", Password: "test", PasswordConfirmation: "test"}

			assert.False(t, r.isValid())
		})

		t.Run("Should fail if the password is empty", func(t *testing.T) {
			r := registerRequest{Email: "test@test.de", Password: "", PasswordConfirmation: "test"}

			assert.False(t, r.isValid())
		})

		t.Run("Should fail if the password confirmation is not equal to the password", func(t *testing.T) {
			r := registerRequest{Email: "test@test.de", Password: "test", PasswordConfirmation: "test1"}

			assert.False(t, r.isValid())
		})

		t.Run("Should succeed if fields are correct", func(t *testing.T) {
			r := registerRequest{Email: "test@test.de", Password: "test", PasswordConfirmation: "test"}

			assert.True(t, r.isValid())
		})
	})
}

func TestLogin(t *testing.T) {
	ctrl := gomock.NewController(t)

	repo := mocks.NewMockRepository(ctrl)
	hasher := mocks.NewMockHasher(ctrl)
	tokenGenerator := mocks.NewMockTokenGenerator(ctrl)
	controller := Controller{repo, hasher, tokenGenerator}

	t.Run("should return 400 BAD REQUEST if payload is not json", func(t *testing.T) {
		body := strings.NewReader(`{"invalid json`)

		// given
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/api/auth/login", body)

		// when
		controller.Login(w, r)

		// then
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should return 400 BAD REQUEST if payload is incomplete", func(t *testing.T) {
		tests := []io.Reader{
			strings.NewReader(`{"email":"test@test.com"}`),
			strings.NewReader(`{"password":"test"}`),
		}

		for _, test := range tests {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/api/auth/login", test)

			// when
			controller.Login(w, r)

			// then
			assert.Equal(t, http.StatusBadRequest, w.Code)
		}
	})

	t.Run("should return 401 UNAUTHORIZED if user not found", func(t *testing.T) {
		// given
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/api/auth/login", strings.NewReader(`{"email":"test@test.com","password":"test"}`))

		repo.
			EXPECT().
			FindByEmail("test@test.com").
			Return([]*model.DbUser{}, nil)

		// when
		controller.Login(w, r)

		// then
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("should return 401 UNAUTHORIZED if password is not correct", func(t *testing.T) {
		// given
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/api/auth/login", strings.NewReader(`{"email":"test@test.com","password":"wrong password"}`))

		repo.
			EXPECT().
			FindByEmail("test@test.com").
			Return([]*model.DbUser{{
				Email:    "test@test.com",
				Password: []byte("hashed password"),
			}}, nil)

		hasher.
			EXPECT().
			Validate([]byte("wrong password"), []byte("hashed password")).
			Return(false)

		// when
		controller.Login(w, r)

		// then
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("should return 200 OK", func(t *testing.T) {
		// given
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/api/auth/login", strings.NewReader(`{"email":"test@test.com","password":"test"}`))

		repo.
			EXPECT().
			FindByEmail("test@test.com").
			Return([]*model.DbUser{{
				Email:    "test@test.com",
				Password: []byte("hashed password"),
			}}, nil)

		hasher.
			EXPECT().
			Validate([]byte("test"), []byte("hashed password")).
			Return(true)

		tokenGenerator.
			EXPECT().
			CreateToken(gomockhelpers.Map(map[string]interface{}{
				"email": "test@test.com",
				"exp":   gomock.Any(),
			})).
			Return("token", nil)

			// when
		controller.Login(w, r)

		// then
		res := w.Result()
		var response map[string]interface{}
		err := json.NewDecoder(res.Body).Decode(&response)

		assert.NoError(t, err)
		assert.Equal(t, "token", response["access_token"])
		assert.Equal(t, "Bearer", response["token_type"])
		assert.Equal(t, float64(3600), response["expires_in"])
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestRegister(t *testing.T) {
	ctrl := gomock.NewController(t)

	repo := mocks.NewMockRepository(ctrl)
	hasher := mocks.NewMockHasher(ctrl)
	tokenGenerator := mocks.NewMockTokenGenerator(ctrl)
	controller := Controller{repo, hasher, tokenGenerator}

	t.Run("should return 400 BAD REQUEST if payload is not json", func(t *testing.T) {
		tests := []io.Reader{
			nil,
			strings.NewReader(`{"invalid json`),
		}

		for _, test := range tests {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/api/auth/register", test)

			// when
			controller.Register(w, r)

			// then
			assert.Equal(t, http.StatusBadRequest, w.Code)
		}
	})

	t.Run("should return 400 BAD REQUEST if payload is incomplete", func(t *testing.T) {
		tests := []io.Reader{
			strings.NewReader(`{}`),
			strings.NewReader(`{"email":"test@test.com"}`),
			strings.NewReader(`{"password":"test"}`),
			strings.NewReader(`{"password_confirmation":"test"}`),
		}

		for _, test := range tests {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/api/auth/register", test)

			// when
			controller.Register(w, r)

			// then
			assert.Equal(t, http.StatusBadRequest, w.Code)
		}
	})

	t.Run("should return 500 INTERNAL SERVER ERROR if search for existing user failed", func(t *testing.T) {
		// given
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/api/auth/register",
			strings.NewReader(`{"email":"test@test.com","password":"test", "password_confirmation": "test"}`))

		repo.
			EXPECT().
			FindByEmail("test@test.com").
			Return(nil, errors.New("could not query database"))

		// when
		controller.Register(w, r)

		// then
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("should return 400 BAD REQUEST if user already exists", func(t *testing.T) {
		// given
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/api/auth/register",
			strings.NewReader(`{"email":"test@test.com","password":"test", "password_confirmation": "test"}`))

		repo.
			EXPECT().
			FindByEmail("test@test.com").
			Return([]*model.DbUser{{}}, nil)

		// when
		controller.Register(w, r)

		// then
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should return 500 INTERNAL SERVER ERROR if hashing password failed", func(t *testing.T) {
		// given
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/api/auth/register", strings.NewReader(`{"email":"test@test.com","password":"test", "password_confirmation": "test"}`))

		repo.
			EXPECT().
			FindByEmail("test@test.com").
			Return([]*model.DbUser{}, nil)

		hasher.
			EXPECT().
			Hash([]byte("test")).
			Return(nil, errors.New("could not hash password"))

		// when
		controller.Register(w, r)

		// then
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("should return 500 INTERNAL SERVER ERROR if user could not be created", func(t *testing.T) {
		// given
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/api/auth/register", strings.NewReader(`{"email":"test@test.com","password":"test", "password_confirmation": "test"}`))

		repo.
			EXPECT().
			FindByEmail("test@test.com").
			Return([]*model.DbUser{}, nil)

		hasher.
			EXPECT().
			Hash([]byte("test")).
			Return([]byte("hashed password"), nil)

		repo.
			EXPECT().
			Create([]*model.DbUser{{
				Email:    "test@test.com",
				Password: []byte("hashed password"),
			}}).
			Return(errors.New("could not create user"))

		// when
		controller.Register(w, r)

		// then
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("should return 200 OK", func(t *testing.T) {
		// given
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/api/auth/register",
			strings.NewReader(`{"email":"test@test.com","password":"test", "password_confirmation":"test"}`))

		repo.
			EXPECT().
			FindByEmail("test@test.com").
			Return([]*model.DbUser{}, nil)

		hasher.
			EXPECT().
			Hash([]byte("test")).
			Return([]byte("hashed password"), nil)

		repo.
			EXPECT().
			Create([]*model.DbUser{{
				Email:    "test@test.com",
				Password: []byte("hashed password"),
			}}).
			Return(nil)

		// when
		controller.Register(w, r)

		// then
		assert.Equal(t, http.StatusOK, w.Code)
	})
}
