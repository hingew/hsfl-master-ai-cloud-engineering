package handler

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

func TestLoginHandler(t *testing.T) {
	ctrl := gomock.NewController(t)

	userRepository := mocks.NewMockRepository(ctrl)
	hasher := mocks.NewMockHasher(ctrl)
	tokenGenerator := mocks.NewMockTokenGenerator(ctrl)
	handler := NewLoginHandler(userRepository, hasher, tokenGenerator)

	t.Run("should return 405 METHOD NOT ALLOWED if method is not POST", func(t *testing.T) {
		// given
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/api/auth/login", nil)

		// when
		handler.ServeHTTP(w, r)

		// then
		assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
	})

	t.Run("should return 400 BAD REQUEST if payload is not json", func(t *testing.T) {
		tests := []io.Reader{
			nil,
			strings.NewReader(`{"invalid json`),
		}

		for _, test := range tests {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/api/auth/login", test)

			// when
			handler.ServeHTTP(w, r)

			// then
			assert.Equal(t, http.StatusBadRequest, w.Code)
		}
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
			handler.ServeHTTP(w, r)

			// then
			assert.Equal(t, http.StatusBadRequest, w.Code)
		}
	})

	t.Run("should return 500 INTERNAL SERVER ERROR if search for user failed", func(t *testing.T) {
		// given
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/api/auth/login", strings.NewReader(`{"email":"test@test.com","password":"test"}`))

		userRepository.
			EXPECT().
			FindByEmail("test@test.com").
			Return(nil, errors.New("could not query database"))

		// when
		handler.ServeHTTP(w, r)

		// then
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("should return 401 UNAUTHORIZED if user not found", func(t *testing.T) {
		// given
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/api/auth/login", strings.NewReader(`{"email":"test@test.com","password":"test"}`))

		userRepository.
			EXPECT().
			FindByEmail("test@test.com").
			Return([]*model.DbUser{}, nil)

		// when
		handler.ServeHTTP(w, r)

		// then
		assert.Equal(t, "Basic realm=Restricted", w.Header().Get("WWW-Authenticate"))
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("should return 401 UNAUTHORIZED if password is not correct", func(t *testing.T) {
		// given
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/api/auth/login", strings.NewReader(`{"email":"test@test.com","password":"wrong password"}`))

		userRepository.
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
		handler.ServeHTTP(w, r)

		// then
		assert.Equal(t, "Basic realm=Restricted", w.Header().Get("WWW-Authenticate"))
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("should return 200 OK", func(t *testing.T) {
		// given
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/api/auth/login", strings.NewReader(`{"email":"test@test.com","password":"test"}`))

		userRepository.
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
		handler.ServeHTTP(w, r)

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
