package controller

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	mock "github.com/hingew/hsfl-master-ai-cloud-engineering/creation-service/_mocks"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/lib/model"
	"go.uber.org/mock/gomock"

	"github.com/stretchr/testify/assert"
)

func TestValidation(t *testing.T) {
	t.Run("should be valid if no params are given and the template has no elements", func(t *testing.T) {
		template := &model.PdfTemplate{
			ID:        1,
			Name:      "Test",
			CreatedAt: time.Date(2023, 10, 9, 1, 1, 1, 1, time.UTC),
			UpdatedAt: time.Date(2023, 10, 9, 3, 1, 1, 1, time.UTC),
			Elements:  []model.Element{},
		}

		params := make(map[string]string)

		assert.True(t, isValid(template, params))
	})

	t.Run("should be invalid a param is missing", func(t *testing.T) {
		template := &model.PdfTemplate{
			ID:        1,
			Name:      "Test",
			CreatedAt: time.Date(2023, 10, 9, 1, 1, 1, 1, time.UTC),
			UpdatedAt: time.Date(2023, 10, 9, 3, 1, 1, 1, time.UTC),
			Elements: []model.Element{
				{
					Type:      "text",
					ValueFrom: "value",
					X:         0,
					Y:         0,
					Width:     2,
					Height:    1,
				},
			},
		}

		params := make(map[string]string)

		assert.False(t, isValid(template, params))
	})

	t.Run("should be valid a params are provided", func(t *testing.T) {
		template := &model.PdfTemplate{
			ID:        1,
			Name:      "Test",
			CreatedAt: time.Date(2023, 10, 9, 1, 1, 1, 1, time.UTC),
			UpdatedAt: time.Date(2023, 10, 9, 3, 1, 1, 1, time.UTC),
			Elements: []model.Element{
				{
					Type:      "text",
					ValueFrom: "value",
					X:         0,
					Y:         0,
					Width:     2,
					Height:    1,
				},
			},
		}

		params := make(map[string]string)
        params["value"] = "Hello world!"

		assert.True(t, isValid(template, params))
	})
}

func TestController(t *testing.T) {
	ctrl := gomock.NewController(t)

	pdf := mock.NewMockPdf(ctrl)
	client := mock.NewMockTemplatingServiceClient(ctrl)

	controller := NewController(pdf, client)

	t.Run("Should fail if the id is not a int", func(t *testing.T) {
		// given

		endpoint := fmt.Sprintf("/api/render/%s", "string")
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", endpoint,
			strings.NewReader(`
            { "value": "whatever" }
                `))

		r = r.WithContext(context.WithValue(r.Context(), "id", "string"))

		// when
		controller.CreatePdf(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Should fail if the templating service in unavialable", func(t *testing.T) {
		// given

		var id uint = 1234
		endpoint := fmt.Sprintf("/api/render/%d", id)
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", endpoint,
			strings.NewReader(`
            { "value": "whatever" }
                `))

		r = r.WithContext(context.WithValue(r.Context(), "id", "1234"))

		client.
			EXPECT().
			FetchTemplate(id).
			Return(nil, errors.New("Templating service unavialable")).
			Times(1)

		// when
		controller.CreatePdf(w, r)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
	t.Run("Should fail when the params are not json", func(t *testing.T) {
		// given

		var id uint = 1234
		endpoint := fmt.Sprintf("/api/render/%d", id)
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", endpoint,
			strings.NewReader(`
            { "value": "whatever" 
                `))

		r = r.WithContext(context.WithValue(r.Context(), "id", "1234"))

		template := &model.PdfTemplate{
			ID:        id,
			Name:      "Test",
			CreatedAt: time.Date(2023, 10, 9, 1, 1, 1, 1, time.UTC),
			UpdatedAt: time.Date(2023, 10, 9, 3, 1, 1, 1, time.UTC),
			Elements: []model.Element{
				{
					Type:   "rect",
					X:      0,
					Y:      0,
					Width:  2,
					Height: 1,
				},
			},
		}

		client.
			EXPECT().
			FetchTemplate(id).
			Return(template, nil).
			Times(1)

		// when
		controller.CreatePdf(w, r)

		//then
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Should if there are not params", func(t *testing.T) {
		// given

		var id uint = 1234
		endpoint := fmt.Sprintf("/api/render/%d", id)
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", endpoint,
			strings.NewReader(`
            { }
                `))

		r = r.WithContext(context.WithValue(r.Context(), "id", "1234"))

		template := &model.PdfTemplate{
			ID:        id,
			Name:      "Test",
			CreatedAt: time.Date(2023, 10, 9, 1, 1, 1, 1, time.UTC),
			UpdatedAt: time.Date(2023, 10, 9, 3, 1, 1, 1, time.UTC),
			Elements:  []model.Element{},
		}

		client.
			EXPECT().
			FetchTemplate(id).
			Return(template, nil).
			Times(1)

		// when
		controller.CreatePdf(w, r)

		//then
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/pdf", w.Header().Get("Content-Type"))
	})

	t.Run("Should fail if a param is missing", func(t *testing.T) {
		// given

		var id uint = 1234
		endpoint := fmt.Sprintf("/api/render/%d", id)
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", endpoint,
			strings.NewReader(`
            { "value": "whatever" }
                `))

		r = r.WithContext(context.WithValue(r.Context(), "id", "1234"))

		template := &model.PdfTemplate{
			ID:        id,
			Name:      "Test",
			CreatedAt: time.Date(2023, 10, 9, 1, 1, 1, 1, time.UTC),
			UpdatedAt: time.Date(2023, 10, 9, 3, 1, 1, 1, time.UTC),
			Elements: []model.Element{
				{
					Type:      "text",
					X:         0,
					Y:         0,
					Width:     2,
					Height:    1,
					Font:      "Roboto",
					FontSize:  12,
					ValueFrom: "text",
				},
			},
		}

		client.
			EXPECT().
			FetchTemplate(id).
			Return(template, nil).
			Times(1)

		// when
		controller.CreatePdf(w, r)

		//then
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Should create a pdf without params", func(t *testing.T) {
		// given

		var id uint = 1234
		endpoint := fmt.Sprintf("/api/render/%d", id)
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", endpoint,
			strings.NewReader(`
            { "value": "whatever" }
                `))

		r = r.WithContext(context.WithValue(r.Context(), "id", "1234"))

		template := &model.PdfTemplate{
			ID:        id,
			Name:      "Test",
			CreatedAt: time.Date(2023, 10, 9, 1, 1, 1, 1, time.UTC),
			UpdatedAt: time.Date(2023, 10, 9, 3, 1, 1, 1, time.UTC),
			Elements: []model.Element{
				{
					Type:   "rect",
					X:      0,
					Y:      0,
					Width:  2,
					Height: 1,
				},
			},
		}

		client.
			EXPECT().
			FetchTemplate(id).
			Return(template, nil).
			Times(1)

		// when
		controller.CreatePdf(w, r)

		//then
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/pdf", w.Header().Get("Content-Type"))
	})

	t.Run("Should create a pdf with params", func(t *testing.T) {
		// given

		var id uint = 1234
		endpoint := fmt.Sprintf("/api/render/%d", id)
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", endpoint,
			strings.NewReader(`
            { "text": "whatever" }
                `))

		r = r.WithContext(context.WithValue(r.Context(), "id", "1234"))

		template := &model.PdfTemplate{
			ID:        id,
			Name:      "Test",
			CreatedAt: time.Date(2023, 10, 9, 1, 1, 1, 1, time.UTC),
			UpdatedAt: time.Date(2023, 10, 9, 3, 1, 1, 1, time.UTC),
			Elements: []model.Element{
				{
					Type:      "text",
					X:         0,
					Y:         0,
					Width:     2,
					Height:    1,
					Font:      "courier",
					FontSize:  12,
					ValueFrom: "text",
				},
			},
		}

		client.
			EXPECT().
			FetchTemplate(id).
			Return(template, nil).
			Times(1)

		// when
		controller.CreatePdf(w, r)

		//then
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/pdf", w.Header().Get("Content-Type"))
	})

}
