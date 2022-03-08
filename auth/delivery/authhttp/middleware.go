package authhttp

import (
	"Test_derictory/auth"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"net/http"
)

type AuthMiddleware struct {
	usecase auth.UseCase
	redDB   *redis.Client
}

func NewAuthMiddleware(usecase auth.UseCase, redDB *redis.Client) gin.HandlerFunc {
	return (&AuthMiddleware{usecase: usecase, redDB: redDB}).Handle
}

func (m *AuthMiddleware) Handle(c *gin.Context) {
	token, err := c.Cookie("Authorization")
	if err != nil {
		c.JSON(http.StatusUnauthorized, "Cookie has not found")
	}

	//c.Header("Authorization", fmt.Sprintf("Bearer %v", token))

	//authHeader := c.GetHeader("Authorization")
	if token == "" {
		newErrorResponse(c, 401, "Необходима авторизация ")
		return
	}

	//headerParts := strings.Split(token, "+")
	//if len(headerParts) != 2 {
	//	c.AbortWithStatus(http.StatusUnauthorized)
	//	return
	//}

	//if headerParts[0] != "Bearer" {
	//c.AbortWithStatus(http.StatusUnauthorized)
	//return
	//}

	td, err := m.usecase.ParseToken(c.Request.Context(), token)
	if err != nil {
		logrus.Info(err)
		status := http.StatusInternalServerError
		if err == auth.ErrInvalidAccessToken {
			status = http.StatusUnauthorized
		}

		c.AbortWithStatus(status)
		return
	}

	userID, err := m.redDB.Get(c.Request.Context(), td.AccessUUID).Result()
	if err != nil {
		return
	}

	c.Set(auth.CtxUserKey, userID)

}
