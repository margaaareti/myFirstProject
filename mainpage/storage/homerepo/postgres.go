package homerepo

import (
	"Test_derictory/models"
	"context"
	"fmt"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

const (
	studentTable = "students"
	noteTable    = "notes"
	userTable    = "users"
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

	idStr := strconv.Itoa(int(userId))

	query := fmt.Sprintf("INSERT INTO %s (name, surname, patronymic,isu_number,added_by,title,description) values ($1,$2,$3,$4,$5,$6,$7) RETURNING id", studentTable)
	row := s.db.QueryRow(ctx, query, student.Name, student.Surname, student.Patronymic, student.IsuNumber, idStr, student.Title, student.Description)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	logrus.Infof("A new entry has been added:  %s %s ", student.Surname, student.Name)

	return id, nil

}

func (s *StudentRepo) PullAllNotice(ctx context.Context) ([]models.Student, error) {

	var notes []models.Student

	query := fmt.Sprintf("SELECT id,name,surname,patronymic, added_by, isu_number,title,description,reg_date FROM %s ", studentTable)
	err := pgxscan.Select(ctx, s.db, &notes, query)
	if err != nil {
		return nil, err
	}
	return notes, nil
}

func (s *StudentRepo) DeleteNotice(ctx context.Context, Id int) error {

	query := fmt.Sprintf(`DELETE FROM %s WHERE id=$1`, studentTable)
	_, err := s.db.Exec(ctx, query, Id)
	if err != nil {
		return err
	}
	return nil
}

func (s *StudentRepo) PullById(ctx context.Context, id uint64) (models.Student, error) {
	var student models.Student
	query := fmt.Sprintf(`SELECT * FROM %s WHERE id = %v`, studentTable, id)
	if err := pgxscan.Get(ctx, s.db, &student, query); err != nil {
		return student, errors.New("That student is not found")
	}

	return student, nil
}

func (s *StudentRepo) UpdateEntry(ctx context.Context, userId, studentId uint64, input models.UpdateStudentInput) error {

	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argId))
		args = append(args, *input.Title)
		argId++
	}

	if input.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argId))
		args = append(args, *input.Description)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf("UPDATE %s SET %s WHERE id= %d", studentTable, setQuery, studentId)

	//args = append(args, studentId)

	_, err := s.db.Exec(ctx, query, args...)
	return err
}
