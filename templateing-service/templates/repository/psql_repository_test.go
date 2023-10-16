package repository

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/hingew/hsfl-master-ai-cloud-engineering/templating-service/templates/model"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestPsqlRepository(t *testing.T) {
	db, dbmock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	repository := PsqlRepository{db}

	t.Run("CreateTemplate", func(t *testing.T) {
		t.Run("should insert a template in the db", func(t *testing.T) {
			// given
			template := model.PdfTemplate{
				PdfName: "Test",
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

			elementJson, _ := json.Marshal(template.Elements)

			dbmock.ExpectExec(`INSERT INTO templates (name, created_at, updated_at, elements) VALUES ($1, $2, $3, $4)`).
				WithArgs("Test", sqlmock.AnyArg(), sqlmock.AnyArg(), elementJson).
				WillReturnResult(nil)

			// when
			err := repository.CreateTemplate(&template)

			// then
			assert.NoError(t, err)
			assert.NoError(t, dbmock.ExpectationsWereMet())
		})
	})

	t.Run("GetAllTemplates", func(t *testing.T) {
		t.Run("should return all templates", func(t *testing.T) {
			// given
			dbmock.ExpectQuery(`SELECT id, name, created_at, updated_at, elements FROM templates`).
				WillReturnRows(sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at", "elements"}).
					AddRow(1, "Template1", "2023-10-12 10:27:10.297308+00", "2023-10-12 10:27:10.297308+00", "[{\"type\":\"rect\",\"x\":0,\"y\":0,\"width\":0,\"height\":0,\"value_from\":\"\",\"font\":\"\",\"size\":0},{\"type\":\"text\",\"x\":0,\"y\":0,\"width\":0,\"height\":0,\"value_from\":\"title\",\"font\":\"JetBrainsMono\",\"size\":18}]").
					AddRow(2, "Template2", "2023-10-12 10:05:46.400412+00", "2023-10-12 10:05:46.400412+00", "[]"))

			// when
			templates, err := repository.GetAllTemplates()

			// then
			assert.NoError(t, err)
			assert.NoError(t, dbmock.ExpectationsWereMet())
			assert.Len(t, templates, 2)
			assert.Equal(t, "Template1", templates[0].PdfName)
			assert.Equal(t, time.Date(2023, 10, 12, 10, 27, 10, 297308, time.Local), templates[0].CreatedAt)
			assert.Equal(t, []model.Element{
				{
					Type:   "rect",
					X:      0,
					Y:      0,
					Width:  2,
					Height: 1,
				},
				{
					Type:      "text",
					X:         0,
					Y:         0,
					Width:     0,
					Height:    0,
					ValueFrom: "title",
					Font:      "JetBrainsMono",
					Size:      18,
				},
			}, templates[0].Elements)
			assert.Equal(t, "Template2", templates[1].PdfName)
		})
	})

	t.Run("GetTemplate", func(t *testing.T) {
		t.Run("should return template by id", func(t *testing.T) {
			// given
			var id int64 = 999

			dbmock.ExpectQuery(`SELECT name, created_at, updated_at, elements FROM templates WHERE id = $1`).
				WithArgs(999).
				WillReturnRows(sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at", "elements"}).
					AddRow(999, "Template1", "2023-10-12 10:27:10.297310+00", "2023-10-13 10:31:10.297308+00", "[{\"type\":\"rect\",\"x\":0,\"y\":0,\"width\":0,\"height\":0,\"value_from\":\"\",\"font\":\"\",\"size\":0},{\"type\":\"text\",\"x\":0,\"y\":0,\"width\":0,\"height\":0,\"value_from\":\"title\",\"font\":\"JetBrainsMono\",\"size\":18}]"))

			// when
			template, err := repository.GetTemplate(id)

			// then
			assert.NoError(t, err)
			assert.NotNil(t, template)
			assert.NoError(t, dbmock.ExpectationsWereMet())
			assert.Equal(t, "Template1", template.PdfName)
			assert.Equal(t, time.Date(2023, 10, 12, 10, 31, 10, 297308, time.Local), template.CreatedAt)
			assert.Equal(t, time.Date(2023, 10, 13, 10, 27, 10, 297310, time.Local), template.UpdatedAt)
			assert.Equal(t, []model.Element{
				{
					Type:   "rect",
					X:      0,
					Y:      0,
					Width:  2,
					Height: 1,
				},
				{
					Type:      "text",
					X:         0,
					Y:         0,
					Width:     0,
					Height:    0,
					ValueFrom: "title",
					Font:      "JetBrainsMono",
					Size:      18,
				},
			}, template.Elements)
		})
	})

	t.Run("DeleteTemplate", func(t *testing.T) {
		t.Run("should delete template in batch", func(t *testing.T) {
			// given
			dbmock.ExpectExec(`DELETE FROM templates WHERE id = $1`).
				WithArgs(1).
				WillReturnResult(sqlmock.NewResult(0, 1))

			// when
			err := repository.DeleteTemplate(1)

			// then
			assert.NoError(t, err)
			assert.NoError(t, dbmock.ExpectationsWereMet())
		})
	})
}
