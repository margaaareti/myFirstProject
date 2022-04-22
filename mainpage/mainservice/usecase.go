package mainservice

import "Test_derictory/mainpage"

type HomeUseCase struct {
	CrudPage mainpage.HomePage
}

func NewHomeUseCase(crudpage mainpage.HomeRepo) *HomeUseCase {
	return &HomeUseCase{CrudPage: crudpage}
}
