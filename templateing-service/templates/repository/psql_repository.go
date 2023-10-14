package repository

import (
	"database/sql"
	"encoding/json"
	"reflect"
	"time"

	"github.com/hingew/hsfl-master-ai-cloud-engineering/lib/database"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/templating-service/templates/model"

	_ "github.com/lib/pq"
)

type PsqlRepository struct {
	db *sql.DB
}

func NewPsqlRepository(config database.Config) (*PsqlRepository, error) {
	db, err := sql.Open("postgres", config.Dsn())
	if err != nil {
		return nil, err
	}

	return &PsqlRepository{db}, nil
}

const createTemplatesTable = `
create table if not exists templates (
	id          serial  	primary key,
	name        text    	not null,
	created_at	timestamptz not null,
	updated_at	timestamptz not null,
	elements	json		not null
)
`

func (repo *PsqlRepository) Setup() error {
	_, err := repo.db.Exec(createTemplatesTable)
	return err
}

func (repo *PsqlRepository) CreateTemplate(template *model.PdfTemplate) error {
	elementJson, err := json.Marshal(template.Elements)
	if err != nil {
		return err
	}

	queryPlaceholder := "INSERT INTO templates (name, created_at, updated_at, elements) VALUES ($1, $2, $3, $4)"
	_, err = repo.db.Exec(queryPlaceholder, template.PdfName, time.Now(), time.Now(), elementJson)
	return err
}

func (repo *PsqlRepository) GetAllTemplates() ([]*model.PdfTemplate, error) {
	rows, err := repo.db.Query("SELECT id, name, created_at, updated_at, elements FROM templates")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var templates []*model.PdfTemplate

	for rows.Next() {
		var template model.PdfTemplate
		var elementJson string

		err := rows.Scan(&template.ID, &template.PdfName, &template.CreatedAt, &template.UpdatedAt, &elementJson)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal([]byte(elementJson), &template.Elements)
		if err != nil {
			return nil, err
		}

		templates = append(templates, &template)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return templates, nil
}

func (repo *PsqlRepository) GetTemplate(id int64) (*model.PdfTemplate, error) {
	var template model.PdfTemplate
	var elementJson string

	template.ID = id

	err := repo.db.QueryRow("SELECT name, created_at, updated_at, elements FROM templates WHERE id = $1", id).Scan(
		&template.PdfName,
		&template.CreatedAt,
		&template.UpdatedAt,
		&elementJson,
	)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(elementJson), &template.Elements)
	if err != nil {
		return nil, err
	}

	return &template, nil
}

func (repo *PsqlRepository) UpdateTemplate(id int64, template *model.PdfTemplate) error {
	dbVersion, err := repo.GetTemplate(id)
	if err != nil {
		return err
	}

	if template.PdfName != "" && dbVersion.PdfName != template.PdfName {
		repo.update_template_name(id, template.PdfName)
	}
	dbElementsAreEmpty := areElementListsEqual(dbVersion.Elements, []model.Element{})
	passedElementsAreEmpty := areElementListsEqual(template.Elements, []model.Element{})
	elementsAreEqual := areElementListsEqual(dbVersion.Elements, template.Elements)
	if (!dbElementsAreEmpty && passedElementsAreEmpty) || !elementsAreEqual {
		repo.update_template_elements(id, template.Elements)
	}

	return err
}

func areElementListsEqual(list1 []model.Element, list2 []model.Element) bool {
	if len(list1) != len(list2) {
		return false
	}

	for i, element1 := range list1 {
		element2 := list2[i]

		if !reflect.DeepEqual(element1, element2) {
			return false
		}
	}

	return true
}

func (repo *PsqlRepository) update_template_name(id int64, name string) error {
	_, err := repo.db.Exec("UPDATE templates SET name = $1, updated_at = $2 WHERE id = $3", name, time.Now(), id)
	return err
}

func (repo *PsqlRepository) update_template_elements(id int64, elements []model.Element) error {
	elementJson, err := json.Marshal(elements)
	if err != nil {
		return err
	}

	_, err = repo.db.Exec("UPDATE templates SET elements = $1, updated_at = $2 WHERE id = $3", elementJson, time.Now(), id)
	return err
}

func (repo *PsqlRepository) DeleteTemplate(id int64) error {
	_, err := repo.db.Exec("DELETE FROM templates WHERE id = $1", id)
	return err
}
