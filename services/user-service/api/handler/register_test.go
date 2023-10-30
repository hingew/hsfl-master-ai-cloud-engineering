package handler

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	mocks "github.com/hingew/hsfl-master-ai-cloud-engineering/user-service/_mocks"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/user-service/user/model"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestRegisterHandler(t *testing.T) {
	ctrl := gomock.NewController(t)

	hasher := mocks.NewMockHasher(ctrl)
	userRepository := mocks.NewMockRepository(ctrl)
	handler := Register(userRepository, hasher)

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
			strings.NewReader(`{}`),
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

	t.Run("should return 500 INTERNAL SERVER ERROR if search for existing user failed", func(t *testing.T) {
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

	t.Run("should return 409 CONFLICT if user already exists", func(t *testing.T) {
		// given
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/api/auth/login", strings.NewReader(`{"email":"test@test.com","password":"test"}`))

		userRepository.
			EXPECT().
			FindByEmail("test@test.com").
			Return([]*model.DbUser{{}}, nil)

		// when
		handler.ServeHTTP(w, r)

		// then
		assert.Equal(t, http.StatusConflict, w.Code)
	})

	t.Run("should return 500 INTERNAL SERVER ERROR if hashing password failed", func(t *testing.T) {
		// given
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/api/auth/login", strings.NewReader(`{"email":"test@test.com","password":"test"}`))

		userRepository.
			EXPECT().
			FindByEmail("test@test.com").
			Return([]*model.DbUser{}, nil)

		hasher.
			EXPECT().
			Hash([]byte("test")).
			Return(nil, errors.New("could not hash password"))

		// when
		handler.ServeHTTP(w, r)

		// then
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("should return 500 INTERNAL SERVER ERROR if user could be created", func(t *testing.T) {
		// given
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/api/auth/login", strings.NewReader(`{"email":"test@test.com","password":"test"}`))

		userRepository.
			EXPECT().
			FindByEmail("test@test.com").
			Return([]*model.DbUser{}, nil)

		hasher.
			EXPECT().
			Hash([]byte("test")).
			Return([]byte("hashed password"), nil)

		userRepository.
			EXPECT().
			Create([]*model.DbUser{{
				Email:    "test@test.com",
				Password: []byte("hashed password"),
			}}).
			Return(errors.New("could not create user"))

		// when
		handler.ServeHTTP(w, r)

		// then
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("should return 200 OK", func(t *testing.T) {
		// given
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/api/auth/login", strings.NewReader(`{"email":"test@test.com","password":"test"}`))

		userRepository.
			EXPECT().
			FindByEmail("test@test.com").
			Return([]*model.DbUser{}, nil)

		hasher.
			EXPECT().
			Hash([]byte("test")).
			Return([]byte("hashed password"), nil)

		userRepository.
			EXPECT().
			Create([]*model.DbUser{{
				Email:    "test@test.com",
				Password: []byte("hashed password"),
			}}).
			Return(nil)

		// when
		handler.ServeHTTP(w, r)

		// then
		assert.Equal(t, http.StatusOK, w.Code)
	})
}
