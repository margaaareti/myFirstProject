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
	useCase auth.UseCase
	redDB   *redis.Client
}

func NewAuthMiddleware(useCase auth.UseCase, redDB *redis.Client) gin.HandlerFunc {
	return (&AuthMiddleware{useCase: useCase, redDB: redDB}).Handle
}

func (m *AuthMiddleware) Handle(c *gin.Context) {

	aToken, err := c.Cookie("AccessToken")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, "Необходима авторизация")
		return
	}

	rToken, err := c.Cookie("RefreshToken")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, "Нет данных авторизации")
		return
	}

	//c.Header("Authorization", fmt.Sprintf("Bearer %v", token))

	//authHeader := c.GetHeader("Authorization")
	if aToken == "" {
		newErrorResponse(c, 401, "Необходима авторизация ")
		return
	}
	if rToken == "" {
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

	ad, err := m.useCase.ParseAcsToken(c.Request.Context(), aToken)
	if err != nil {

		uuid, username, err := m.useCase.ParseRefToken(c.Request.Context(), rToken)
		if err != nil {
			newErrorResponse(c, 401, err.Error())
			return
		}

		deletedId, delErr := m.useCase.DeleteTokens(c.Request.Context(), uuid)
		if delErr != nil || deletedId == 0 {
			newErrorResponse(c, 401, delErr.Error())
			return
		}
		tokens, id, err := m.useCase.CreateTokens(c.Request.Context(), username, deletedId)

		saveErr := m.useCase.CreateAuth(c.Request.Context(), id, tokens)
		if saveErr != nil {
			newErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		c.SetCookie("AccessToken", tokens.AccessToken, 60*60*24, "/", "localhost", false, true)
		c.SetCookie("RefreshToken", tokens.RefreshToken, 60*60*24, "/", "localhost", false, true)

		c.Set(auth.CtxUserKey, id)
		logrus.Info(err)

	} else if err == auth.ErrInvalidAccessToken {
		status := http.StatusInternalServerError
		status = http.StatusUnauthorized
		c.AbortWithStatus(status)
		return

	} else {
		userID, err := m.redDB.Get(c.Request.Context(), ad.AccessUUID).Result()
		if err != nil {
			newErrorResponse(c, 401, err.Error())
			return
		}

		c.Set(auth.CtxUserKey, userID)
	}

}
