package mainservice

import (
	"Test_derictory/mainpage"
	"Test_derictory/models"
	"context"
)

type HomeUseCase struct {
	CrudPage mainpage.HomeRepo
}

func NewHomeUseCase(crudpage mainpage.HomeRepo) *HomeUseCase {
	return &HomeUseCase{CrudPage: crudpage}
}

func (h *HomeUseCase) AddStudent(ctx context.Context, userId uint64, student models.Student) (uint64, error) {
	return h.CrudPage.CreateStudent(ctx, userId, student)
}

func (h *HomeUseCase) GetAllNotice(ctx context.Context) ([]models.Student, error) {
	return h.CrudPage.PullAllNotice(ctx)
}

func (h *HomeUseCase) GetById(ctx context.Context, Id uint64) (models.Student, error) {
	return h.CrudPage.PullById(ctx, Id)
}

func (h *HomeUseCase) DeleteNoticeByID(ctx context.Context, Id int) error {
	return h.CrudPage.DeleteNotice(ctx, Id)
}

func (h *HomeUseCase) UpdateEntryUseCase(ctx context.Context, userId, studentId uint64, input models.UpdateStudentInput) error {
	if err := input.Validate(); err != nil {
		return err
	}
	return h.CrudPage.UpdateEntry(ctx, userId, studentId, input)
}
