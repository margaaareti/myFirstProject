package authhttp

import (
	"Test_derictory/auth"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

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

	if aToken == "" {
		newErrorResponse(c, 401, "Необходима авторизация ")
		return
	}
	if rToken == "" {
		newErrorResponse(c, 401, "Необходима авторизация ")
		return
	}

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

func GetUserId(c *gin.Context) (uint64, error) {
	id, ok := c.Get(auth.CtxUserKey)
	logrus.Info(id)
	if !ok {
		return 0, errors.New("user id not found")
	}
	//Приводим id в соответствубщему типу
	idStr, ok := id.(string)
	if !ok {
		return 0, errors.New("user id not found")
	}
	u64, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return 0, errors.New("some shit happened")
	}

	return u64, nil

}
