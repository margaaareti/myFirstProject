package authhttp

import (
	"Test_derictory/auth"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"net/http"
)

type tokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type AuthMiddleware struct {
	usecase auth.UseCase
	redDB   *redis.Client
}

func NewAuthMiddleware(usecase auth.UseCase, redDB *redis.Client) gin.HandlerFunc {
	return (&AuthMiddleware{usecase: usecase, redDB: redDB}).Handle
}

func (m *AuthMiddleware) Handle(c *gin.Context) {

	tokens := &tokenResponse{}
	var err error

	tokens.AccessToken, err = c.Cookie("Authorization")
	if err != nil {
		c.JSON(http.StatusUnauthorized, "Cookie has not found")
	}

	tokens.RefreshToken, err = c.Cookie("RefreshToken")
	if err != nil {
		c.JSON(http.StatusUnauthorized, "Cookie has not found")
	}

	//c.Header("Authorization", fmt.Sprintf("Bearer %v", token))

	//authHeader := c.GetHeader("Authorization")
	if tokens.AccessToken == "" {
		newErrorResponse(c, 401, "Необходима авторизация ")
		return
	}
	if tokens.RefreshToken == "" {
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

	td, err := m.usecase.ParseToken(c.Request.Context(), tokens.AccessToken)
	if err != nil {
		td, userID, err := m.usecase.ParseRefresh(c.Request.Context(), tokens.RefreshToken)
		tokens.AccessToken = td.AccessToken
		tokens.RefreshToken = td.RefreshToken
		c.Set(auth.CtxUserKey, userID)
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
