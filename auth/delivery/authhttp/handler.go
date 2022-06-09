package authhttp

import (
	"Test_derictory/auth"
	"Test_derictory/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

const MAX_UPLOAD_SIZE = 10 << 20

var IMAGE_TYPES = map[string]interface{}{
	"image/jpeg": nil,
	"image/png":  nil,
}

type Handler struct {
	useCase auth.UseCase
}

func NewHandler(useCase auth.UseCase) *Handler {
	return &Handler{useCase: useCase}
}

type signInput struct {
	Username string `form:"username"`
	Password string `form:"password"`
}

func (h *Handler) SignUp(c *gin.Context) {

	if c.Request.Method != "POST" {
		newErrorResponse(c, http.StatusMethodNotAllowed, "ForbiddenMethod")
		return
	}

	var input models.User2

	if err := c.Bind(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	_, err := h.useCase.SignUp(c.Request.Context(), input)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.Redirect(303, "/auth/cong")

}

func (h *Handler) SignIn(c *gin.Context) {

	if c.Request.Method != "POST" {
		newErrorResponse(c, http.StatusMethodNotAllowed, "ForbiddenMethod")
	}

	var inp signInput

	if err := c.Bind(&inp); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.useCase.SignIn(c.Request.Context(), inp.Username, inp.Password)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	token, data, err := h.useCase.CreateTokens(c.Request.Context(), user)

	saveErr := h.useCase.CreateAuth(c.Request.Context(), data.Id, token)
	if saveErr != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.SetCookie("AccessToken", token.AccessToken, 60*60*24, "/", "localhost", false, true)
	c.SetCookie("RefreshToken", token.RefreshToken, 60*60*24, "/", "localhost", false, true)

	c.JSON(http.StatusOK, map[string]interface{}{
		"access_token":  token.AccessToken,
		"refresh_token": token.RefreshToken,
	})

	c.Redirect(303, "/api/home")

}

func (h *Handler) UploadFile(c *gin.Context) {

	err := c.Request.ParseMultipartForm(MAX_UPLOAD_SIZE) //10mb
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	file, fileHeader, err := c.Request.FormFile("myFile")
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Error Retrieving file from form-data:%v", err.Error()))
		return
	}

	defer file.Close()

	buffer := make([]byte, fileHeader.Size)
	file.Read(buffer)
	fileType := http.DetectContentType(buffer)
	fmt.Println(fileType)

	//Validate file type
	if _, ex := IMAGE_TYPES[fileType]; !ex {
		newErrorResponse(c, http.StatusBadRequest, "Неверный тип файла")
		return
	}

	status := fmt.Sprintf("File has been uploaded: %+v", fileHeader.Filename)
	c.JSON(http.StatusOK, map[string]interface{}{
		"status": status,
		"size":   fileHeader.Size,
		"name":   fileHeader.Header,
	})

}
