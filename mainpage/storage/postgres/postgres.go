package postgres

import (
	"Test_derictory/models"
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	studentTable = "students"
	noteTable    = "notes"
)

type StudentRepo struct {
	db *pgxpool.Pool
}

func NewAuthPostgres(db *pgxpool.Pool) *StudentRepo {
	return &StudentRepo{db: db}
}

func (s *StudentRepo) CreateStudent(ctx context.Context, userId uint64, student models.Student) (uint64, error) {

}
