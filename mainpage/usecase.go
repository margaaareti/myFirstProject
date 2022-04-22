package mainpage

import (
	"Test_derictory/models"
	"context"
)

type HomePage interface {
	AddStudent(ctx context.Context, userId int, student models.Student) (int, error)
}
