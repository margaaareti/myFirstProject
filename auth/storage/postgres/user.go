package postgres

import (
	"Test_derictory/auth/storage/email"
	"Test_derictory/models"
	"Test_derictory/server/repository"
	"context"
	"fmt"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type AuthPostgres struct {
	db *pgxpool.Pool
}

func NewAuthPostgres(db *pgxpool.Pool) *AuthPostgres {
	return &AuthPostgres{db: db}
}

//Create user in DB
func (r *AuthPostgres) CreateUser(ctx context.Context, user models.User2) (uint64, error) {

	var id uint64

	var mailStatus int
	var userStatus int
	isExistGet := fmt.Sprintf(`SELECT COUNT(id), (SELECT COUNT(id) FROM %[1]s WHERE username = $2) FROM %[1]s WHERE email = $1`, repository.UserTable)
	isExistRow := r.db.QueryRow(ctx, isExistGet, user.Email, user.Username)
	if err := isExistRow.Scan(&mailStatus, &userStatus); err != nil {
		return 0, err
	}

	if userStatus != 0 {
		return 0, errors.New(userAlrExist)
	} else if mailStatus != 0 {
		return 0, errors.New(emailAlrExist)
	} else {
		query := fmt.Sprintf("INSERT INTO %s (name,username,password,email) values($1,$2,$3,$4) RETURNING id", repository.UserTable)
		row := r.db.QueryRow(ctx, query, user.Name, user.Username, user.Password, user.Email)
		if err := row.Scan(&id); err != nil {
			return 0, err
		}
		if err := email.SendEmail(user.Email); err != nil {
			return 0, err
		}
		logrus.Infof("A new user has been registered:  %s ", user.Username)
		return id, nil

	}

}

func (r *AuthPostgres) GetUser(ctx context.Context, username, password string) (*models.User2, error) {

	newUser := new(models.User2)

	query := fmt.Sprintf(`SELECT Id,name,surname,patronymic FROM %s WHERE username = '%s' AND password = '%s'`, repository.UserTable, username, password)
	if err := pgxscan.Get(ctx, r.db, newUser, query); err != nil {
		return nil, errors.New("That user is not found")
	}

	return newUser, nil

}
