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

func TestController(t *testing.T) {
	ctrl := gomock.NewController(t)

	pdf := mock.NewMockPdf(ctrl)
	client := mock.NewMockTemplatingServiceClient(ctrl)

	controller := NewController(pdf, client)

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
			GetTemplate(id).
			Return(nil, errors.New("Templating service unavialable")).
			Times(1)

		// when
		controller.CreatePdf(w, r)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("Should create a pdf", func(t *testing.T) {
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
			GetTemplate(id).
			Return(template, nil).
			Times(1)

		// when
		controller.CreatePdf(w, r)

		//then
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/pdf", w.Header().Get("Content-Type"))
	})
}
