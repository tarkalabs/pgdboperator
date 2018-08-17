package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/teris-io/shortid"
)

type PostgresHandler struct {
	dsn string
	db  *sql.DB
}

func NewPostgresHandler(dsn string) *PostgresHandler {
	return &PostgresHandler{dsn: dsn}
}

func (pgh *PostgresHandler) IsConnected() bool {
	return pgh.db != nil
}

func (pgh *PostgresHandler) EstablishConnection() error {
	conn, err := sql.Open("postgres", pgh.dsn)
	if err != nil {
		return err
	}
	pgh.db = conn
	return nil
}

func (pgh *PostgresHandler) CreateRandomRole() (string, error) {
	if !pgh.IsConnected() {
		return "", fmt.Errorf("Handler is currently not connected to a database")
	}
	roleName, err := shortid.Generate()
	passwd := GeneratePassword(32)
	if err != nil {
		return "", fmt.Errorf("Unable to generate id %+v", err)
	}
	stmt := fmt.Sprintf("CREATE ROLE \"%s\" LOGIN PASSWORD '%s'", roleName, passwd)
	fmt.Println(stmt)
	_, err = pgh.db.Exec(stmt)
	if err != nil {
		return "", fmt.Errorf("Unable to create role %s - %+v", roleName, err)
	}
	return roleName, nil
}

func (pgh *PostgresHandler) CreateDatabase(dbName string, role string) error {
	if !pgh.IsConnected() {
		return fmt.Errorf("Handler is currently not connected to a database")
	}
	tx, err := pgh.db.Begin()
	if err != nil {
		return fmt.Errorf("Unable to create transaction while creating database %s for role %s - %+v", dbName, role, err)
	}
	stmt := fmt.Sprintf("CREATE DATABASE \"%s\" OWNER \"%s\"", dbName, role)
	_, err = tx.Exec(stmt)
	if err != nil {
		defer tx.Rollback()
		return fmt.Errorf("Unable to create database %s for role %s - %+v", dbName, role, err)
	}
	stmt = fmt.Sprintf("REVOKE ALL PRIVILEGES ON DATABASE \"%s\" FROM PUBLIC; GRANT ALL PRIVILEGES ON DATABASE \"%s\" TO \"%s\";", dbName, dbName, role)
	_, err = tx.Exec(stmt)
	if err != nil {
		defer tx.Rollback()
		return fmt.Errorf("Unable to set permissions for database %s for role %s - %+v", dbName, role, err)
	}
	return tx.Commit()
}
