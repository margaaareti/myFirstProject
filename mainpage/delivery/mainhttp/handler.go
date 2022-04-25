package mainhttp

import (
	"Test_derictory/auth"
	"Test_derictory/auth/delivery/authhttp"
	"Test_derictory/mainpage"
	"Test_derictory/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

type HomeHandler struct {
	handHome mainpage.HomePage
	auth     auth.UseCase
}

func NewHomeHandler(handHome mainpage.HomePage, auth auth.UseCase) *HomeHandler {
	return &HomeHandler{handHome: handHome,
		auth: auth}
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

func (h *HomeHandler) LogOut(c *gin.Context) {

	if c.Request.Method != "POST" {
		newErrorResponse(c, http.StatusMethodNotAllowed, "ForbiddenMethod")
	}

	aToken, err := c.Cookie("AccessToken")
	rToken, err := c.Cookie("RefreshToken")

	myIn, err := h.auth.ParseAcsToken(c.Request.Context(), aToken)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	myIn.RefreshUUID, _, err = h.auth.ParseRefToken(c.Request.Context(), rToken)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
	}

	deleted, delErr := h.auth.LogOut(c.Request.Context(), myIn.AccessUUID, myIn.RefreshUUID)
	if delErr != nil && deleted == 0 {
		c.JSON(401, "unauthorized")
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"Status": myIn,
	})

	c.Redirect(303, "/auth/sign-in")

}

func (h *HomeHandler) CreateEntry(c *gin.Context) {

	userId, err := authhttp.GetUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	var input models.Student
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.handHome.AddStudent(c.Request.Context(), userId, input)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (h *HomeHandler) AllNotes(c *gin.Context) {

	entries, err := h.handHome.GetAllNotice(c.Request.Context())
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"data": entries,
	})
}
