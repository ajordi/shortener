package mysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"github.com/ajordi/shortener/shortener"
)

type mysqlRepository struct {
	client *sql.DB
}

func newMysqlClient(connectionString string) (*sql.DB, error) {
	client, err := sql.Open("mysql", connectionString)
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}

	return client, nil
}

func NewMysqlRepository(serverName string, user string, password string, dbName string) (shortener.RedirectRepository, error) {
	connectionString := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=true&multiStatements=true", user, password, serverName, dbName)
	client, err := newMysqlClient(connectionString)
	if err != nil {
		return nil, errors.Wrap(err, "repository.NewMysqlRepository")
	}

	repo := &mysqlRepository{}
	repo.client = client

	return repo, nil
}

func (m *mysqlRepository) generateKey(code string) string {
	return fmt.Sprintf("redirect:%s", code)
}

func (m *mysqlRepository) Find(code string) (*shortener.Redirect, error) {
	r := &shortener.Redirect{}

	err := m.client.QueryRow("select code, url, createdAt from redirects WHERE code = ?", code).Scan(&r.Code, &r.URL, &r.CreatedAt)

	if err != nil {
		fmt.Println(err.Error())
	}

	return r, nil
}

func (m *mysqlRepository) Store(redirect *shortener.Redirect) error {

	stmtIns, err := m.client.Prepare("INSERT INTO redirects(code, url, createdAt) VALUES(?,?,?)")
	if err != nil {
		panic(err.Error())
	}

	_, err = stmtIns.Exec(redirect.Code, redirect.URL, redirect.CreatedAt)

	return err
}

