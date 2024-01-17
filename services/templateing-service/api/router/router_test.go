package router

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	mock_middleware "github.com/hingew/hsfl-master-ai-cloud-engineering/lib/_mocks"
	mock_controller "github.com/hingew/hsfl-master-ai-cloud-engineering/templateing-service/_mocks"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
)

func TestRouter(t *testing.T) {
	ctrl := gomock.NewController(t)

	templatingController := mock_controller.NewMockControllerInterface(ctrl)
	authMiddleware := mock_middleware.NewMockAuthMiddleInterface(ctrl)
	router := NewTemplateRouter(templatingController, authMiddleware)

	passAuthMiddleware := func(w http.ResponseWriter, r *http.Request) {
		authMiddleware.
			EXPECT().
			AuthMiddleware(w, r, gomock.Any()).
			Do(func(w http.ResponseWriter, r *http.Request, next func(r *http.Request)) {
				next(r)
			}).
			Times(1)
	}

	t.Run("/templates", func(t *testing.T) {
		t.Run("should return 404 NOT FOUND if method is not GET or POST", func(t *testing.T) {
			tests := []string{"DELETE", "PUT", "HEAD", "CONNECT", "OPTIONS", "TRACE", "PATCH"}

			for _, test := range tests {
				// given
				w := httptest.NewRecorder()
				r := httptest.NewRequest(test, "/api/templates", nil)

				passAuthMiddleware(w, r)

				// when
				router.ServeHTTP(w, r)

				// then
				assert.Equal(t, http.StatusNotFound, w.Code)
			}
		})

		t.Run("should call GET handler", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/templates", nil)

			templatingController.
				EXPECT().
				GetAllTemplates(w, r).
				Times(1)

			passAuthMiddleware(w, r)

			// when
			router.ServeHTTP(w, r)

			// then
			assert.Equal(t, http.StatusOK, w.Code)
		})

		t.Run("should call POST handler", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/api/templates", nil)

			templatingController.
				EXPECT().
				CreateTemplate(w, r).
				Times(1)

			passAuthMiddleware(w, r)

			// when
			router.ServeHTTP(w, r)

			// then
			assert.Equal(t, http.StatusOK, w.Code)
		})
	})

	t.Run("/templates/:id", func(t *testing.T) {
		t.Run("should return 404 NOT FOUND if method is not GET, DELETE or PUT", func(t *testing.T) {
			tests := []string{"POST", "HEAD", "CONNECT", "OPTIONS", "TRACE", "PATCH"}

			for _, test := range tests {
				// given
				w := httptest.NewRecorder()
				r := httptest.NewRequest(test, "/api/templates/1", nil)

				passAuthMiddleware(w, r)
				// when
				router.ServeHTTP(w, r)

				// then
				assert.Equal(t, http.StatusNotFound, w.Code)
			}
		})

		t.Run("should call GET handler", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/templates/1", nil)

			templatingController.
				EXPECT().
				GetTemplate(w, r.WithContext(context.WithValue(r.Context(), "id", "1"))).
				Times(1)

			passAuthMiddleware(w, r)

			// when
			router.ServeHTTP(w, r)

			// then
			assert.Equal(t, http.StatusOK, w.Code)
		})

		t.Run("should call PUT handler", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("PUT", "/api/templates/1", nil)

			templatingController.
				EXPECT().
				UpdateTemplate(w, r.WithContext(context.WithValue(r.Context(), "id", "1"))).
				Times(1)

			passAuthMiddleware(w, r)

			// when
			router.ServeHTTP(w, r)

			// then
			assert.Equal(t, http.StatusOK, w.Code)
		})

		t.Run("should call DELETE handler", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("DELETE", "/api/templates/1", nil)

			templatingController.
				EXPECT().
				DeleteTemplate(w, r.WithContext(context.WithValue(r.Context(), "id", "1"))).
				Times(1)

			passAuthMiddleware(w, r)

			// when
			router.ServeHTTP(w, r)

			// then
			assert.Equal(t, http.StatusOK, w.Code)
		})
	})
}
