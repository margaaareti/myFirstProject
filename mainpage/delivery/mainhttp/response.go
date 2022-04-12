package mainhttp

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func newErrorResponse(c *gin.Context, statusCode int, message string) {
	logrus.Errorf(message)
	c.AbortWithStatusJSON(statusCode, message)
}
