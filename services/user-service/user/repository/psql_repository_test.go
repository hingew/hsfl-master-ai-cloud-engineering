package repository

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/user-service/user/model"
	"github.com/stretchr/testify/assert"
)

func TestPsqlRepository(t *testing.T) {
	db, dbmock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	repository := PsqlRepository{db}

	t.Run("Create", func(t *testing.T) {
		t.Run("should return error if executing query failed", func(t *testing.T) {
			// given
			users := []*model.DbUser{{
				Email:    "doesnt matter",
				Password: []byte("doesnt matter"),
			}}

			dbmock.
				ExpectExec(`insert into users`).
				WillReturnError(errors.New("database error"))

			// when
			err := repository.Create(users)

			// then
			assert.Error(t, err)
		})

		t.Run("should insert users in batches", func(t *testing.T) {
			// given
			users := []*model.DbUser{
				{
					Email:    "test@test.com",
					Password: []byte("test"),
				},
				{
					Email:    "abc@abc.com",
					Password: []byte("abc"),
				},
			}

			dbmock.
				ExpectExec(`insert into users \(email, password\) values \(\$1,\$2\),\(\$3,\$4\)`).
				WithArgs("test@test.com", []byte("test"), "abc@abc.com", []byte("abc")).
				WillReturnResult(sqlmock.NewResult(0, 2))

			// when
			err := repository.Create(users)

			// then
			assert.NoError(t, err)
		})
	})

	t.Run("FindByEmail", func(t *testing.T) {
		t.Run("should return error if executing query failed", func(t *testing.T) {
			// given
			email := "test@test.com"

			dbmock.
				ExpectQuery(`select email, password from users where email = \$1`).
				WillReturnError(errors.New("database error"))

			// when
			users, err := repository.FindByEmail(email)

			// then
			assert.Error(t, err)
			assert.Nil(t, users)
		})

		t.Run("should return users by email", func(t *testing.T) {
			// given
			email := "test@test.com"

			dbmock.
				ExpectQuery(`select email, password from users where email = \$1`).
				WillReturnRows(sqlmock.NewRows([]string{"email", "password"}).AddRow("test@test.com", []byte("hash")))

			// when
			users, err := repository.FindByEmail(email)

			// then
			assert.NoError(t, err)
			assert.Len(t, users, 1)
		})
	})

	t.Run("Delete", func(t *testing.T) {
		t.Run("should return error if executing query failed", func(t *testing.T) {
			// given
			users := []*model.DbUser{
				{
					Email:    "test@test.com",
					Password: []byte("hash"),
				},
				{
					Email:    "abc@abc.com",
					Password: []byte("hash"),
				},
			}

			dbmock.
				ExpectExec(`delete from users`).
				WillReturnError(errors.New("database error"))

			// when
			err := repository.Delete(users)

			// then
			assert.Error(t, err)
		})

		t.Run("should delete users in batches", func(t *testing.T) {
			// given
			users := []*model.DbUser{
				{
					Email:    "test@test.com",
					Password: []byte("hash"),
				},
				{
					Email:    "abc@abc.com",
					Password: []byte("hash"),
				},
			}

			dbmock.
				ExpectExec(`delete from users where email in \(\$1,\$2\)`).
				WithArgs("test@test.com", "abc@abc.com").
				WillReturnResult(sqlmock.NewResult(0, 2))

			// when
			err := repository.Delete(users)

			// then
			assert.NoError(t, err)
		})
	})
}
