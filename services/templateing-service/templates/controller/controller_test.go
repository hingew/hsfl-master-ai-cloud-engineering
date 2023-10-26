package controller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	mock_repository "github.com/hingew/hsfl-master-ai-cloud-engineering/templating-service/_mock"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/templating-service/templates/model"
	"github.com/stretchr/testify/assert"
)

func TestController(t *testing.T) {
	ctrl := gomock.NewController(t)

	productRepository := mock_repository.NewMockIRepository(ctrl)
	controller := ControllerImp{productRepository}

	t.Run("GetAllTemplates", func(t *testing.T) {
		t.Run("should return 500 INTERNAL SERVER ERROR if query failed", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/templates", nil)

			productRepository.
				EXPECT().
				GetAllTemplates().
				Return(nil, errors.New("query failed")).
				Times(1)

			// when
			controller.GetAllTemplates(w, r)

			// then
			assert.Equal(t, http.StatusInternalServerError, w.Code)
		})

		t.Run("should return all templates", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/templates", nil)

			productRepository.
				EXPECT().
				GetAllTemplates().
				Return(&[]model.PdfTemplate{
					{
						ID:        999,
						Name:      "Test",
						CreatedAt: time.Date(2023, 10, 9, 1, 1, 1, 1, time.UTC),
						UpdatedAt: time.Date(2023, 10, 9, 3, 1, 1, 1, time.UTC),
						Elements:  []model.Element{},
					},
				}, nil).
				Times(1)

			// when
			controller.GetAllTemplates(w, r)

			// then
			res := w.Result()
			var response []model.PdfTemplate
			err := json.NewDecoder(res.Body).Decode(&response)

			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, w.Code)
			assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
			assert.Len(t, response, 1)
			assert.Equal(t, uint(999), response[0].ID)
			assert.Equal(t, string("Test"), response[0].Name)
			assert.Equal(t, time.Date(2023, 10, 9, 1, 1, 1, 1, time.UTC), response[0].CreatedAt)
			assert.Equal(t, time.Date(2023, 10, 9, 3, 1, 1, 1, time.UTC), response[0].UpdatedAt)
			assert.Len(t, response[0].Elements, 0)
		})
	})

	t.Run("GetTemplate", func(t *testing.T) {
		t.Run("should return all products", func(t *testing.T) {
			// given
			var id uint = 1234
			endpoint := fmt.Sprintf("/templates/%d", id)
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", endpoint, nil)
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
			productRepository.
				EXPECT().
				GetTemplateById(id).
				Return(template, nil).
				Times(1)

			// when
			controller.GetTemplate(w, r)

			// then
			res := w.Result()
			var response model.PdfTemplate
			err := json.NewDecoder(res.Body).Decode(&response)

			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, w.Code)
			assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
			assert.Equal(t, uint(id), response.ID)
			assert.Equal(t, string("Test"), response.Name)
			assert.Equal(t, time.Date(2023, 10, 9, 1, 1, 1, 1, time.UTC), response.CreatedAt)
			assert.Equal(t, time.Date(2023, 10, 9, 3, 1, 1, 1, time.UTC), response.UpdatedAt)
			assert.Len(t, response.Elements, 1)
			assert.Equal(t, string("rect"), response.Elements[0].Type)
			assert.Equal(t, int(0), response.Elements[0].X)
			assert.Equal(t, int(0), response.Elements[0].Y)
			assert.Equal(t, int(2), response.Elements[0].Width)
			assert.Equal(t, int(1), response.Elements[0].Height)
		})

		t.Run("should return 500 INTERNAL SERVER ERROR if query failed", func(t *testing.T) {
			// given
			var id uint = 1234
			endpoint := fmt.Sprintf("/templates/%d", id)
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", endpoint, nil)
			r = r.WithContext(context.WithValue(r.Context(), "id", "1234"))

			productRepository.
				EXPECT().
				GetTemplateById(id).
				Return(nil, errors.New("query failed")).
				Times(1)

			// when
			controller.GetTemplate(w, r)

			// then
			assert.Equal(t, http.StatusInternalServerError, w.Code)
		})

		t.Run("should return 400 BAD REQUEST if ID is not numeric", func(t *testing.T) {
			// given
			id := "13$8ho169"
			endpoint := fmt.Sprintf("/templates/%s", id)
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", endpoint, nil)
			r = r.WithContext(context.WithValue(r.Context(), "id", "13$8ho169"))

			// when
			controller.GetTemplate(w, r)

			// then
			assert.Equal(t, http.StatusBadRequest, w.Code)
		})
	})

	t.Run("CreateTemplate", func(t *testing.T) {
		t.Run("should return 400 BAD REQUEST if payload is not json", func(t *testing.T) {
			tests := []io.Reader{
				nil,
				strings.NewReader(`{"invalid`),
			}

			for _, test := range tests {
				// given
				w := httptest.NewRecorder()
				r := httptest.NewRequest("POST", "/templates", test)

				// when
				controller.CreateTemplate(w, r)

				// then
				assert.Equal(t, http.StatusBadRequest, w.Code)
			}
		})

		t.Run("should return 500 INTERNAL SERVER ERROR if persisting failed", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/templates",
				strings.NewReader(`{
					"name": "My table template",
					"elements": [
						{ 
							"type": "rect",
							"x": 0, 
							"y": 0, 
							"width": 0, 
							"height": 0
						},
						{
							"type": "text",
							"x": 0, 
							"y": 0, 
							"width": 0, 
							"height": 0,
							"value_from": "title",
							"font": "JetBrainsMono",
							"size": 18
						}
					]
				}`))

			productRepository.
				EXPECT().
				CreateTemplate(gomock.Any()).
				Return(nil, errors.New("database error"))

			// when
			controller.CreateTemplate(w, r)

			// then
			assert.Equal(t, http.StatusInternalServerError, w.Code)
		})

		t.Run("should create new template", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/templates",
				strings.NewReader(`{
					"name": "My table template",
					"elements": [
						{ 
							"type": "rect",
							"x": 0, 
							"y": 0, 
							"width": 0, 
							"height": 0
						},
						{
							"type": "text",
							"x": 0, 
							"y": 0, 
							"width": 0, 
							"height": 0,
							"value_from": "title",
							"font": "JetBrainsMono",
							"size": 18
						}
					]
				}`))

			id := uint(1)

			productRepository.
				EXPECT().
				CreateTemplate(gomock.Any()).
				Return(&id, nil)

			// when
			controller.CreateTemplate(w, r)

			// then
			assert.Equal(t, http.StatusOK, w.Code)
		})
	})

	t.Run("UpdateTemplate", func(t *testing.T) {
		t.Run("should return 400 BAD REQUEST if template id is not numerical", func(t *testing.T) {
			// given
			id := "13$8ho169"
			endpoint := fmt.Sprintf("/templates/%s", id)
			w := httptest.NewRecorder()
			r := httptest.NewRequest("PUT", endpoint, nil)
			r = r.WithContext(context.WithValue(r.Context(), "id", "13$8ho169"))

			// when
			controller.UpdateTemplate(w, r)

			// then
			assert.Equal(t, http.StatusBadRequest, w.Code)
		})

		t.Run("should return 400 BAD REQUEST if payload is not json", func(t *testing.T) {
			tests := []io.Reader{
				nil,
				strings.NewReader(`{"invalid`),
			}

			for _, test := range tests {
				// given
				w := httptest.NewRecorder()
				r := httptest.NewRequest("PUT", "/templates/1", test)
				r = r.WithContext(context.WithValue(r.Context(), "id", "1"))

				// when
				controller.UpdateTemplate(w, r)

				// then
				assert.Equal(t, http.StatusBadRequest, w.Code)
			}
		})

		t.Run("should return 500 INTERNAL SERVER ERROR if query failed", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("PUT", "/templates/1",
				strings.NewReader(`{"name": "New Name"}`))
			r = r.WithContext(context.WithValue(r.Context(), "id", "1"))

			request := model.PdfTemplate{
				Name: "New Name",
			}

			productRepository.
				EXPECT().
				UpdateTemplate(uint(1), request).
				Return(errors.New("database error"))

			// when
			controller.UpdateTemplate(w, r)

			// then
			assert.Equal(t, http.StatusInternalServerError, w.Code)
		})

		t.Run("should update one product", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("PUT", "/templates/1",
				strings.NewReader(`{"name": "New Name"}`))
			r = r.WithContext(context.WithValue(r.Context(), "id", "1"))

			request := model.PdfTemplate{
				Name: "New Name",
			}

			productRepository.
				EXPECT().
				UpdateTemplate(uint(1), request).
				Return(nil)

			// when
			controller.UpdateTemplate(w, r)

			// then
			assert.Equal(t, http.StatusOK, w.Code)
		})
	})

	t.Run("DeleteProduct", func(t *testing.T) {
		t.Run("should return 400 BAD REQUEST if product id is not numerical", func(t *testing.T) {
			// given
			id := "13$8ho169"
			endpoint := fmt.Sprintf("/templates/%s", id)
			w := httptest.NewRecorder()
			r := httptest.NewRequest("DELETE", endpoint, nil)
			r = r.WithContext(context.WithValue(r.Context(), "id", "13$8ho169"))

			// when
			controller.DeleteTemplate(w, r)

			// then
			assert.Equal(t, http.StatusBadRequest, w.Code)
		})

		t.Run("should return 500 INTERNAL SERVER ERROR if query fails", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("DELETE", "/templates/1", nil)
			r = r.WithContext(context.WithValue(r.Context(), "id", "1"))

			productRepository.
				EXPECT().
				DeleteTemplate(uint(1)).
				Return(errors.New("database error"))

			// when
			controller.DeleteTemplate(w, r)

			// then
			assert.Equal(t, http.StatusInternalServerError, w.Code)
		})

		t.Run("should return 200 OK", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("DELETE", "/templates/1", nil)
			r = r.WithContext(context.WithValue(r.Context(), "id", "1"))

			productRepository.
				EXPECT().
				DeleteTemplate(uint(1)).
				Return(nil)

			// when
			controller.DeleteTemplate(w, r)

			// then
			assert.Equal(t, http.StatusOK, w.Code)
		})
	})
}
