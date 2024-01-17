package repository

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/hingew/hsfl-master-ai-cloud-engineering/lib/containerhelpers"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/lib/database"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/lib/model"
	"github.com/stretchr/testify/assert"
)

func TestIntegrationGormPsqlRepository(t *testing.T) {
	postgres, err := containerhelpers.StartPostgres()
	if err != nil {
		t.Fatalf("could not start postgres container: %s", err.Error())
	}

	t.Cleanup(func() {
		postgres.Terminate(context.Background())
	})

	port, err := postgres.MappedPort(context.Background(), "5432")
	if err != nil {
		t.Fatalf("could not get database container port: %s", err.Error())
	}

	repository, err := NewGormPsqlRepository(database.PsqlConfig{
		Host:     "0.0.0.0",
		Port:     port.Int(),
		Username: "postgres",
		Password: "postgres",
		Database: "postgres",
	})

	if err != nil {
		t.Fatalf("could not create user repository: %s", err.Error())
	}

	sqlDb, err := repository.db.DB()
	if err != nil {
		t.Fatalf("could not access underlying sql.DB of gorm.DB: %s", err.Error())
	}

	t.Cleanup(clearTables(t, sqlDb, []string{"elements", "pdf_templates"}))

	t.Run("CreateTemplate", func(t *testing.T) {
		// arrange
		if err := repository.Migrate(); err != nil {
			t.Fatalf("could not migrate repo: %s", err.Error())
		}
		t.Cleanup(clearTables(t, sqlDb, []string{"elements", "pdf_templates"}))

		template := model.PdfTemplate{Name: "test", Elements: []model.Element{
			{Type: "rect", X: 0, Y: 0, Width: 2, Height: 1},
			{Type: "circle", X: 3, Y: 4, Width: 3, Height: 3},
		}}

		// act
		id, err := repository.CreateTemplate(template)

		// assert
		numberOfTemplates := numberOfRows(t, sqlDb, "pdf_templates")
		numberOfElements := numberOfRows(t, sqlDb, "elements")

		assert.NoError(t, err)
		assert.NotNil(t, id)
		assert.Equal(t, 1, numberOfTemplates)
		assert.Equal(t, 2, numberOfElements)
	})

	t.Run("GetTemplateById", func(t *testing.T) {
		// arrange
		if err := repository.Migrate(); err != nil {
			t.Fatalf("could not migrate repo: %s", err.Error())
		}

		t.Cleanup(clearTables(t, sqlDb, []string{"elements", "pdf_templates"}))

		template := model.PdfTemplate{Name: "test", Elements: []model.Element{
			{Type: "rect", X: 0, Y: 0, Width: 2, Height: 1},
			{Type: "circle", X: 3, Y: 4, Width: 3, Height: 3},
		}}
		insertTemplate(t, sqlDb, &template)
		assert.Equal(t, 1, numberOfRows(t, sqlDb, "pdf_templates"))
		assert.Equal(t, 2, numberOfRows(t, sqlDb, "elements"))

		// act
		dbTemplate, err := repository.GetTemplateById(template.ID)

		//assert
		assert.NoError(t, err)
		assert.NotNil(t, dbTemplate)
		assert.Truef(t, templatesAreEqual(t, template, *dbTemplate), "Expected:\t%s\nActual:\t\t%s", template, *dbTemplate)
	})

	t.Run("GetAllTemplates", func(t *testing.T) {
		// arrange
		if err := repository.Migrate(); err != nil {
			t.Fatalf("could not migrate repo: %s", err.Error())
		}

		t.Cleanup(clearTables(t, sqlDb, []string{"elements", "pdf_templates"}))

		template1 := model.PdfTemplate{Name: "test", Elements: []model.Element{
			{Type: "rect", X: 0, Y: 0, Width: 2, Height: 1},
			{Type: "circle", X: 3, Y: 4, Width: 3, Height: 3},
		}}
		template2 := model.PdfTemplate{Name: "test2", Elements: []model.Element{}}
		insertTemplate(t, sqlDb, &template1)
		insertTemplate(t, sqlDb, &template2)
		assert.Equal(t, 2, numberOfRows(t, sqlDb, "pdf_templates"))
		assert.Equal(t, 2, numberOfRows(t, sqlDb, "elements"))

		// act
		dbTemplates, err := repository.GetAllTemplates()

		//assert
		assert.NoError(t, err)
		assert.NotNil(t, dbTemplates)
		assert.Truef(t, templatesAreEqual(t, template1, (*dbTemplates)[0]), "Expected:\t%s\nActual:\t\t%s", template1, (*dbTemplates)[0])
		assert.Truef(t, templatesAreEqual(t, template2, (*dbTemplates)[1]), "Expected:\t%s\nActual:\t\t%s", template1, (*dbTemplates)[1])
	})

	t.Run("GetAllTemplates", func(t *testing.T) {
		// arrange
        if err := repository.Migrate(); err != nil {
            t.Fatalf("could not migrate repo: %s", err.Error())
        }
		t.Cleanup(clearTables(t, sqlDb, []string{"elements", "pdf_templates"}))

		template := model.PdfTemplate{Name: "test", Elements: []model.Element{
			{Type: "rect", X: 0, Y: 0, Width: 2, Height: 1},
			{Type: "circle", X: 3, Y: 4, Width: 3, Height: 3},
		}}
		insertTemplate(t, sqlDb, &template)
		assert.Equal(t, 1, numberOfRows(t, sqlDb, "pdf_templates"))
		assert.Equal(t, 2, numberOfRows(t, sqlDb, "elements"))

		// act
		err := repository.DeleteTemplate(template.ID)

		//assert
		assert.NoError(t, err)
		assert.Equal(t, 0, numberOfRows(t, sqlDb, "pdf_templates"))
		assert.Equal(t, 0, numberOfRows(t, sqlDb, "elements"))
	})

	t.Run("UpdateTemplate", func(t *testing.T) {
		// arrange
        if err := repository.Migrate(); err != nil {
            t.Fatalf("could not migrate repo: %s", err.Error())
        }
		t.Cleanup(clearTables(t, sqlDb, []string{"elements", "pdf_templates"}))

		template := model.PdfTemplate{Name: "test", Elements: []model.Element{
			{Type: "rect", X: 0, Y: 0, Width: 2, Height: 1},
			{Type: "circle", X: 3, Y: 4, Width: 3, Height: 3},
		}}
		insertTemplate(t, sqlDb, &template)
		assert.Equal(t, 1, numberOfRows(t, sqlDb, "pdf_templates"))
		assert.Equal(t, 2, numberOfRows(t, sqlDb, "elements"))

		// act
		newTemplate := copyPdfTemplate(template)
		newTemplate.Name = "new name"
		newTemplate.Elements = []model.Element{
			{Type: "rect", X: 2, Y: 1, Width: 60, Height: 0},
		}
		err := repository.UpdateTemplate(newTemplate.ID, newTemplate)

		//assert
		assert.NoError(t, err)

		dbTemplate, err := repository.GetTemplateById(template.ID)
		assert.NoError(t, err)
		assert.NotNil(t, dbTemplate)
		assert.Equal(t, newTemplate.ID, dbTemplate.ID)
		assert.Equal(t, newTemplate.Name, dbTemplate.Name)
		assert.Equal(t, newTemplate.CreatedAt, dbTemplate.CreatedAt)
		assert.NotEqual(t, newTemplate.UpdatedAt, dbTemplate.UpdatedAt)
		assert.Equal(t, 1, numberOfRows(t, sqlDb, "elements"))
		assert.Falsef(t, elementsAreEqual(t, template.Elements[0], *&dbTemplate.Elements[0]), "Expected:\t%s\nActual:\t\t%s", template.Elements[0], *&dbTemplate.Elements[0])
	})
}

func clearTables(t *testing.T, db *sql.DB, tabelNames []string) func() {
	return func() {
		for _, tabelName := range tabelNames {
			sqlCmd := fmt.Sprintf("delete from %s", tabelName)

			if _, err := db.Exec(sqlCmd); err != nil {
				t.Logf("could not delete rows from %s: %s", tabelName, err.Error())
				t.FailNow()
			}
		}
	}
}

func assertTableExists(t *testing.T, db *sql.DB, name string, columns []string) {
	rows, err := db.Query(`select column_name from information_schema.columns where table_name = $1`, name)
	if err != nil {
		t.Fail()
		return
	}

	scannedCols := make(map[string]struct{})
	for rows.Next() {
		var column string
		if err := rows.Scan(&column); err != nil {
			t.Logf("expected")
			t.FailNow()
		}

		scannedCols[column] = struct{}{}
	}

	if len(scannedCols) == 0 {
		t.Logf("expected table '%s' to exist, but not found", name)
		t.FailNow()
	}

	for _, col := range columns {
		if _, ok := scannedCols[col]; !ok {
			t.Logf("expected table '%s' to have column '%s'", name, col)
			t.Fail()
		}
	}
}

func getTemplateFromDatabase(t *testing.T, db *sql.DB, template model.PdfTemplate) *model.PdfTemplate {
	row := db.QueryRow(`select id, updated_at, created_at, name  from pdf_templates where name = $1`, template.Name)

	var dbTemplate model.PdfTemplate
	if err := row.Scan(&dbTemplate.ID, &dbTemplate.UpdatedAt, &dbTemplate.CreatedAt, &dbTemplate.Name); err != nil {
		return nil
	}

	for _, element := range template.Elements {
		var dbElement = getElementFromDatabase(t, db, element)
		if dbElement != nil {
			dbTemplate.Elements = append(dbTemplate.Elements, *dbElement)
		}
	}

	return &dbTemplate
}

func getElementFromDatabase(t *testing.T, db *sql.DB, element model.Element) *model.Element {
	row := db.QueryRow(`select id, type, x, y, width, height, value_from, font, font_size, pdf_template_id from elements where type = $1 and x = $1 and y = $2`, element.Type, element.X, element.Y)

	var dbElement model.Element
	if err := row.Scan(&element.ID, &element.Type, &element.X, &element.Y, &element.Width, &element.Height, &element.ValueFrom, &element.Font, &element.FontSize, &element.PdfTemplateID); err != nil {
		return nil
	}

	t.Log("getElementFromDatabase: ", dbElement)

	return &dbElement
}

func insertTemplate(t *testing.T, db *sql.DB, template *model.PdfTemplate) {
	var id uint
	var createdAt time.Time
	var updatedAt time.Time

	err := db.QueryRow("INSERT INTO pdf_templates (created_at, updated_at, name) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at", time.Now(), time.Now(), template.Name).Scan(&id, &createdAt, &updatedAt)
	if err != nil {
		t.Logf("Fehler beim Einfügen in pdf_templates: %s", err.Error())
		t.FailNow()
	}
	template.ID = id
	template.CreatedAt = createdAt
	template.UpdatedAt = updatedAt

	var elementId uint
	var pdfTemplateId uint
	for i, element := range template.Elements {
		err := db.QueryRow(`INSERT INTO elements (type, x, y, width, height, value_from, font, font_size, pdf_template_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id, pdf_template_id`, element.Type, element.X, element.Y, element.Width, element.Height, element.ValueFrom, element.Font, element.FontSize, template.ID).Scan(&elementId, &pdfTemplateId)
		if err != nil {
			t.Logf("Fehler beim Einfügen von Elementen: %s", err.Error())
			t.FailNow()
		}
		template.Elements[i].ID = elementId
		template.Elements[i].PdfTemplateID = pdfTemplateId
	}
}

func elementsAreEqual(t *testing.T, a model.Element, b model.Element) bool {
	return a.ID == b.ID && a.Type == b.Type && a.X == b.X && a.Y == b.Y && a.Font == b.Font && a.FontSize == b.FontSize && a.Height == b.Height && a.Width == b.Width && a.ValueFrom == b.ValueFrom && a.PdfTemplateID == b.PdfTemplateID
}

func templatesAreEqual(t *testing.T, a model.PdfTemplate, b model.PdfTemplate) bool {
	baseDataIsEqual := a.ID == b.ID && a.CreatedAt == b.CreatedAt && a.UpdatedAt == b.UpdatedAt && a.Name == b.Name
	equalElements := len(a.Elements) == len(b.Elements)
	if !equalElements {
		return false
	}
	for i, elementA := range a.Elements {
		equalElements = elementsAreEqual(t, elementA, b.Elements[i])

		if !equalElements {
			return false
		}
	}

	return baseDataIsEqual && equalElements
}

func numberOfRows(t *testing.T, db *sql.DB, tableName string) int {
	var count int
	query := "SELECT COUNT(*) FROM " + tableName

	err := db.QueryRow(query).Scan(&count)
	if err != nil {
		t.Logf("Could not evaluate number of rows in %s: %s", tableName, err)
		t.FailNow()
	}

	return count
}

func copyPdfTemplate(template model.PdfTemplate) model.PdfTemplate {
	copy := model.PdfTemplate{
		ID:        template.ID,
		CreatedAt: template.CreatedAt,
		UpdatedAt: template.UpdatedAt,
		Name:      template.Name,
	}

	copy.Elements = make([]model.Element, len(template.Elements))
	for i, element := range template.Elements {
		copy.Elements[i] = model.Element{
			ID:            element.ID,
			Type:          element.Type,
			X:             element.X,
			Y:             element.Y,
			Width:         element.Width,
			Height:        element.Height,
			ValueFrom:     element.ValueFrom,
			Font:          element.Font,
			FontSize:      element.FontSize,
			PdfTemplateID: element.PdfTemplateID,
		}
	}

	return copy
}
