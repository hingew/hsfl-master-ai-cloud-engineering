package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	mocks "github.com/hingew/hsfl-master-ai-cloud-engineering/user-service/_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestRouter(t *testing.T) {
	ctrl := gomock.NewController(t)

	registerHandler := mocks.NewMockHandler(ctrl)
	loginHandler := mocks.NewMockHandler(ctrl)
	router := New(registerHandler, loginHandler)

	t.Run("should run register handler", func(t *testing.T) {
		// given
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/api/auth/register", nil)

		registerHandler.
			EXPECT().
			ServeHTTP(w, r).
			Times(1)

		// when
		router.ServeHTTP(w, r)

		// then
		assert.Equal(t, http.StatusOK, w.Code)
		assert.True(t, ctrl.Satisfied())
	})

	t.Run("should run login handler", func(t *testing.T) {
		// given
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/api/auth/login", nil)

		loginHandler.
			EXPECT().
			ServeHTTP(w, r).
			Times(1)

		// when
		router.ServeHTTP(w, r)

		// then
		assert.Equal(t, http.StatusOK, w.Code)
		assert.True(t, ctrl.Satisfied())
	})

	t.Run("should return 404 NOT FOUND if target is unknown", func(t *testing.T) {
		// given
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/api/auth/unknown", nil)

		// when
		router.ServeHTTP(w, r)

		// then
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}
