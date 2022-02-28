//create a new DB
package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	"time"
)

//Represent our tables like a constants

const UserTable = "users"

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func NewPostgresDB(cfg Config) (*pgxpool.Pool, error) {

	db, err := pgxpool.Connect(context.Background(), fmt.Sprintf("host =%s port=%s user=%s dbname=%s password=%s sslmode=%s", cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode))
	if err != nil {
		return nil, errors.Errorf("Unable to connect to database: %s", err)
	}

	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = db.Ping(context.Background())
	if err != nil {
		return nil, err
	}
	return db, nil
}
