package authhttp

import (
	"Test_derictory/auth"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

type AuthMiddleware struct {
	usecase auth.UseCase
	redDB   *redis.Client
}

func NewAuthMiddleware(usecase auth.UseCase, redDB *redis.Client) gin.HandlerFunc {
	return (&AuthMiddleware{usecase: usecase, redDB: redDB}).Handle
}

func (m *AuthMiddleware) Handle(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		newErrorResponse(c, 401, "Необходима авторизация")
		return
	}

	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if headerParts[0] != "Bearer" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	user, err := m.usecase.ParseToken(c.Request.Context(), headerParts[1])
	if err != nil {
		logrus.Info(err)
		status := http.StatusInternalServerError
		if err == auth.ErrInvalidAccessToken {
			status = http.StatusUnauthorized
		}

		c.AbortWithStatus(status)
		return
	}

	userID, err := m.redDB.Get(c.Request.Context(), user.AccessUUID).Result()
	if err != nil {
		return
	}

	c.Set(auth.CtxUserKey, userID)

}
