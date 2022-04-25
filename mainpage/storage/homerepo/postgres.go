package homerepo

import (
	"Test_derictory/models"
	"context"
	"fmt"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

const (
	studentTable = "students"
	noteTable    = "notes"
)

type Execute struct {
	Name    string `json:"name" db:"name"`
	Surname string `json:"surname" db:"surname"`
}

type StudentRepo struct {
	db *pgxpool.Pool
}

func NewStudentRepo(db *pgxpool.Pool) *StudentRepo {
	return &StudentRepo{db: db}
}

func (s *StudentRepo) CreateStudent(ctx context.Context, userId uint64, student models.Student) (uint64, error) {

	var id uint64

	query := fmt.Sprintf("INSERT INTO %s (name, surname, patronymic,isu_number,added_by,title,description) values ($1,$2,$3,$4,$5,$6,$7) RETURNING id", studentTable)
	row := s.db.QueryRow(ctx, query, student.Name, student.Surname, student.Patronymic, student.IsuNumber, userId, student.Title, student.Description)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	logrus.Infof("A new entry has been added:  %s %s ", student.Surname, student.Name)

	return id, nil

}

func (s *StudentRepo) PullAllNotice(ctx context.Context) ([]models.Student, error) {

	var notes []models.Student

	query := fmt.Sprintf("SELECT id,name,surname,patronymic,isu_number,added_by FROM %s", studentTable)
	err := pgxscan.Select(ctx, s.db, &notes, query)
	if err != nil {
		return nil, err
	}
	return notes, nil
}
