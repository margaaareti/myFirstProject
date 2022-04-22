package mainhttp

import (
	"Test_derictory/auth"
	"Test_derictory/mainpage"
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

	/*authHeader := c.GetHeader("Authorization")

	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 {
		newErrorResponse(c, http.StatusUnauthorized, "Чтобы выйти- сперва следует зайти")
		return
	}

	if headerParts[0] != "Bearer" {
		newErrorResponse(c, http.StatusUnauthorized, "Чтобы выйти- сперва следует зайти")
		return
	}*/

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

}
