package mainhttp

import (
	"Test_derictory/auth"
	"Test_derictory/auth/delivery/authhttp"
	"Test_derictory/mainpage"
	"Test_derictory/models"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
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
	UserId, ok := c.Get(auth.CtxUserId)
	if !ok {
		newErrorResponse(c, 401, "Необходима авторизация")
		return
	}
	UserName, ok := c.Get(auth.CtxUserName)
	if !ok {
		newErrorResponse(c, 401, "Необходима авторизация")
		return
	}
	UserSur, ok := c.Get(auth.CtxUserSurname)
	if !ok {
		newErrorResponse(c, 401, "Необходима авторизация")
		return
	}

	//UserName, ok2 := c.Get(auth.CtxUserName)
	//if !ok2 {
	//	newErrorResponse(c, 401, "Необходима авторизация aaaaaaaaaaaaa")
	//	return
	//}

	c.HTML(http.StatusOK, "home.html", gin.H{
		"Id":      UserId,
		"Name":    UserName,
		"Surname": UserSur,
	})

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
		logrus.Info("1")
		logrus.Infof("user id is %s", userId)
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	var input models.Student
	if len(input.IsuNumber) != 6 {
		newErrorResponse(c, http.StatusBadRequest, "Неккоректный номер ису")
		return
	}
	if err := c.BindJSON(&input); err != nil {
		logrus.Info("2")
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.handHome.AddStudent(c.Request.Context(), userId, input)
	if err != nil {
		logrus.Info("33")
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (h *HomeHandler) GetAllNotes(c *gin.Context) {

	entries, err := h.handHome.GetAllNotice(c.Request.Context())
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"data": entries,
	})
}

func (h *HomeHandler) GetById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	logrus.Infof("id is %v", id)

	var form models.Student

	form, err = h.handHome.GetById(c.Request.Context(), uint64(id))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"data": form,
	})
}

func (h *HomeHandler) DeleteNoteById(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	err = h.handHome.DeleteNoticeByID(c.Request.Context(), id)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"data": "Ok",
	})
}

func (h *HomeHandler) UpdateItem(c *gin.Context) {
	userId, err := authhttp.GetUserId(c) //Функция определена в middleware.go
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	var input models.UpdateStudentInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.handHome.UpdateEntryUseCase(c.Request.Context(), userId, uint64(id), input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"Status": "Ok",
	})

}

/*func (h *HomeHandler) UploadFile(c *gin.Context) {

err := c.Request.ParseMultipartForm(10 << 20) //10mb
if err != nil {
	newErrorResponse(c, http.StatusBadRequest, err.Error())
	return
}

file, err := c.FormFile("myFile")
if err != nil {
	newErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Error Retrieving file from form-data:%v", err.Error()))
	return
}

status := fmt.Sprintf("File has been uploaded:%+v\n", file.Filename)
c.JSON(http.StatusOK, map[string]interface{}{
	"status":     status,
	"size":       file.Size,
	"mimeHeader": file.Header,
})
*/
