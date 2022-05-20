package mainpage

import (
	"Test_derictory/models"
	"context"
)

type HomeRepo interface {
	CreateStudent(ctx context.Context, userId uint64, student models.Student) (uint64, error)
	PullAllNotice(ctx context.Context) ([]models.Student, error)
	DeleteNotice(ctx context.Context, Id int) error
}
