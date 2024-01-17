package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/hingew/hsfl-master-ai-cloud-engineering/lib/database"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/user-service/user/model"

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

const createUsersTable = `
create table if not exists users (
	email    varchar(100) not null unique,
	password bytea        not null,
	primary key (email)
)
`

func (repo *PsqlRepository) Migrate() error {
	_, err := repo.db.Exec(createUsersTable)
	return err
}

const createUsersBatchQuery = `
insert into users (email, password) values %s
`

func (repo *PsqlRepository) Create(users []*model.DbUser) error {
	placeholders := make([]string, len(users))
	values := make([]interface{}, len(users)*2)

	for i := 0; i < len(users); i++ {
		placeholders[i] = fmt.Sprintf("($%d,$%d)", i*2+1, i*2+2)
		values[i*2+0] = users[i].Email
		values[i*2+1] = users[i].Password
	}

	query := fmt.Sprintf(createUsersBatchQuery, strings.Join(placeholders, ","))
	_, err := repo.db.Exec(query, values...)
	return err
}

const findUsersByEmailQuery = `
select email, password from users where email = $1
`

func (repo *PsqlRepository) FindByEmail(email string) ([]*model.DbUser, error) {
	rows, err := repo.db.Query(findUsersByEmailQuery, email)
	if err != nil {
		return nil, err
	}

	var users []*model.DbUser
	for rows.Next() {
		user := model.DbUser{}
		if err := rows.Scan(&user.Email, &user.Password); err != nil {
			return nil, err
		}

		users = append(users, &user)
	}

	return users, nil
}

const deleteUsersBatchQuery = `
delete from users where email in (%s)
`

func (repo *PsqlRepository) Delete(users []*model.DbUser) error {
	placeholders := make([]string, len(users))
	emails := make([]interface{}, len(users))

	for i := 0; i < len(users); i++ {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		emails[i] = users[i].Email
	}

	query := fmt.Sprintf(deleteUsersBatchQuery, strings.Join(placeholders, ","))
	_, err := repo.db.Exec(query, emails...)
	return err
}
