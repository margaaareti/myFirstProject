package homehttp

import (
	"Test_derictory/auth"
	"Test_derictory/crudpage"
	"github.com/gin-gonic/gin"
	"net/http"
)

type HomeHandler struct {
	handHome crudpage.HomeUsecase
	auth     auth.UseCase
}

func NewCrdHandler(handHome crudpage.HomeUsecase) *HomeHandler {
	return &HomeHandler{handHome: handHome}
}

func (h *HomeHandler) ShowPage(c *gin.Context) {
	User, ok := c.Get(auth.CtxUserKey)
	if !ok {
		newErrorResponse(c, 401, "Необходима авторизация")
	} else {
		c.HTML(http.StatusOK, "home.html", gin.H{
			"Name": User,
		})

	}
}
