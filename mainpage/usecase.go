package mainpage

import (
	"Test_derictory/models"
	"context"
)

type HomePage interface {
	AddStudent(ctx context.Context, userId uint64, student models.Student) (uint64, error)
	GetAllNotice(ctx context.Context) ([]models.Student, error)
	DeleteNoticeByID(ctx context.Context, Id int) error
}
