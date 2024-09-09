package postgres

import (
	"auth_service/config"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func ConnectDB(cnf config.DatabaseConfig) (*sql.DB, error) {
	dns := fmt.Sprintf("user=%s host=%s port=%s password=%s dbname=%s sslmode=disable", cnf.User, cnf.Host, cnf.Port, cnf.Password, cnf.DBName)

	db, err := sql.Open("postgres", dns)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
