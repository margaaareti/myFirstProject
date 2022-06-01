package mainpage

import (
	"Test_derictory/models"
	"context"
)

type HomeRepo interface {
	CreateStudent(ctx context.Context, userId uint64, student models.Student) (uint64, error)
	PullAllNotice(ctx context.Context) ([]models.Student, error)
	DeleteNotice(ctx context.Context, Id int) error
	PullById(ctx context.Context, id uint64) (models.Student, error)
	UpdateEntry(ctx context.Context, userId, studentId uint64, input models.UpdateStudentInput) error
}
