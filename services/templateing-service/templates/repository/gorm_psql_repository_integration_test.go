package repository

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/hingew/hsfl-master-ai-cloud-engineering/lib/containerhelpers"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/lib/database"
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
	fmt.Print(repository)
	//t.Cleanup(clearTables(t, repository.db, []string{"pdf_templates", "elements"}))
}

func clearTables(t *testing.T, db *sql.DB, tabelNames []string) func() {
	return func() {
		if _, err := db.Exec("delete from users"); err != nil {
			for _, tabelName := range tabelNames {
				sqlCmd := fmt.Sprintf("delete from ", tabelName)

				if _, err := db.Exec(sqlCmd); err != nil {
					t.Logf("could not delete rows from %s: %s", tabelName, err.Error())
					t.FailNow()
				}
			}
		}
	}
}
