package homehttp

import (
	"Test_derictory/crudpage"
	"github.com/gin-gonic/gin"
)

func RegisterHTTPEndPoints(c *gin.RouterGroup, crd crudpage.HomeUsecase) {
	cr := NewCrdHandler(crd)

	c.GET("/home", cr.ShowPage)

}
