package router

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRouter(t *testing.T) {
	t.Run("should return 404 NOT FOUND if the path is unknown", func(t *testing.T) {
		// given
		router := New()
		router.GET("/route/without/params", func(w http.ResponseWriter, r *http.Request) {})

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/foobar", nil)

		//when
		router.ServeHTTP(w, r)

		//then
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("should match the correct handler", func(t *testing.T) {
		// given
		router := New()
		router.GET("/foobar", func(w http.ResponseWriter, r *http.Request) {})

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/foobar", nil)

		//when
		router.ServeHTTP(w, r)

		//then
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("should match the correct handler with params", func(t *testing.T) {
		//given
		router := New()
		var ctx context.Context
		router.GET("/foobar/:id/has/:params", func(w http.ResponseWriter, r *http.Request) {
			ctx = r.Context()
		})

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/foobar/7/has/params", nil)

		//when
		router.ServeHTTP(w, r)

		//then
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "7", ctx.Value("id"))
		assert.Equal(t, "params", ctx.Value("params"))
	})

	t.Run("should always match the correct method", func(t *testing.T) {
		//given
		router := New()

		router.GET("/get", func(w http.ResponseWriter, r *http.Request) {})
		router.POST("/post", func(w http.ResponseWriter, r *http.Request) {})
		router.PUT("/put", func(w http.ResponseWriter, r *http.Request) {})
		router.DELETE("/delete", func(w http.ResponseWriter, r *http.Request) {})

		// For get
		for _, method := range []string{http.MethodPost, http.MethodPut, http.MethodDelete} {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(method, "/get", nil)

			//when
			router.ServeHTTP(w, r)

			//then
			assert.Equal(t, http.StatusNotFound, w.Code)
		}

		// for post
		for _, method := range []string{http.MethodGet, http.MethodPut, http.MethodDelete} {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(method, "/post", nil)

			//when
			router.ServeHTTP(w, r)

			//then
			assert.Equal(t, http.StatusNotFound, w.Code)
		}

		// for put
		for _, method := range []string{http.MethodGet, http.MethodPost, http.MethodDelete} {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(method, "/put", nil)

			//when
			router.ServeHTTP(w, r)

			//then
			assert.Equal(t, http.StatusNotFound, w.Code)
		}

		// for delete
		for _, method := range []string{http.MethodGet, http.MethodPost, http.MethodPut} {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(method, "/delete", nil)

			//when
			router.ServeHTTP(w, r)

			//then
			assert.Equal(t, http.StatusNotFound, w.Code)
		}
	})

	t.Run("should not be possible to declare two routes with same route and path", func(t *testing.T) {
		// given
		assert.Panics(t, func() {
			router := New()
			router.GET("/foobar", func(w http.ResponseWriter, r *http.Request) {})
			router.GET("/foobar", func(w http.ResponseWriter, r *http.Request) {})

		})
	})

}
