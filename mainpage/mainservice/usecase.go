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
