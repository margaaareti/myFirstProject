package authhttp

import (
	"Test_derictory/auth"
	"github.com/gin-gonic/gin"
	"net/http"
)

func RegisterHTTPEndPoints(router *gin.Engine, uc auth.UseCase) {
	h := NewHandler(uc)

	authEndpoints := router.Group("/auth")
	{
		authEndpoints.GET("/sign-up", func(c *gin.Context) {
			c.HTML(http.StatusOK, "registration_form.html", nil)
		})

		authEndpoints.POST("/sign-up", h.SignUp)

		authEndpoints.GET("/cong", func(c *gin.Context) {
			c.HTML(http.StatusOK, "cong_form.html", nil)
		})

		authEndpoints.GET("/sign-in", func(c *gin.Context) {
			c.HTML(http.StatusOK, "signInForm.html", nil)
		})

		authEndpoints.POST("/sign-in", h.SignIn)
	}
}
