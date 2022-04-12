package mainhttp

import (
	"Test_derictory/auth"
	"Test_derictory/mainpage"
	"github.com/gin-gonic/gin"
)

func RegisterHTTPEndPoints(c *gin.RouterGroup, crd mainpage.HomePage, auth auth.UseCase) {
	cr := NewHomeHandler(crd, auth)

	c.GET("/home", cr.ShowPage)
	c.POST("/log-out", cr.LogOut)

}
