package service

import (
	"Test_derictory/auth"
	"Test_derictory/models"
	"context"
	"crypto/sha1"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

const (
	salt       = "bcb545454yh5p5HG"
	tokenTTL   = 1 * time.Minute
	refreshTTL = 24 * 7 * time.Hour
	signingKey = "QHhpZGlF2DG3SD3F3G2SDF3H4vCg=="
	refreshKey = "ds7B989umHJ98opi;m2"
)

type AuthClaims struct {
	jwt.StandardClaims
	Username   string    `json:"username"`
	Name       string    `json:"name"`
	Surname    string    `json:"surname"`
	Patronymic string    `json:"patronymic"`
	UserID     uint64    `json:"userID"`
	TokenUUID  uuid.UUID `json:"access_uuid"`
}

type NewTokens struct {
	AccessToken  string
	RefreshToken string
}

type AuthUseCase struct {
	repo auth.UserRepository
	stg  auth.TokenStorage
}

func NewAuthUseCase(repo auth.UserRepository, stg auth.TokenStorage) *AuthUseCase {
	return &AuthUseCase{repo: repo,
		stg: stg}

}

func (a *AuthUseCase) SignUp(ctx context.Context, user models.User2) (uint64, error) {
	user.Password = GeneratePasswordHash(user.Password)

	if err := isEmailValid(user.Email); err != nil {
		return 0, err

	}

	return a.repo.CreateUser(ctx, user)
}

//Hash for password
func GeneratePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

func (a *AuthUseCase) SignIn(ctx context.Context, username, password string) (*models.User2, error) {
	passwordHash := GeneratePasswordHash(password)

	user, err := a.repo.GetUser(ctx, username, passwordHash)
	if err != nil {
		return nil, err
	}

	return user, nil

}

func (a *AuthUseCase) CreateTokens(ctx context.Context, user *models.User2) (*models.TokenDetails, *models.User2, error) {

	var err error
	//Generate Access Token
	td := &models.TokenDetails{}
	td.AtExpires = time.Now().Add(tokenTTL).Unix()
	td.AccessUuid, err = uuid.NewV4()
	if err != nil {
		return nil, nil, err
	}

	td.RtExpires = time.Now().Add(refreshTTL).Unix()
	td.RefreshUuid, err = uuid.NewV4()
	if err != nil {
		return nil, nil, err
	}

	atClaims := &AuthClaims{
		UserID:     user.Id,
		Username:   user.Username,
		Name:       user.Name,
		Surname:    user.Surname,
		Patronymic: user.Patronymic,
		TokenUUID:  td.AccessUuid,
		//AccessUUID: td.AccessUuid,
		StandardClaims: jwt.StandardClaims{
			// Токен перестанет быть валидным через 15 минут с момента его генерации
			ExpiresAt: td.AtExpires,
			// Время генерации токена
			IssuedAt: time.Now().Unix(),
		},
	}

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(signingKey))
	if err != nil {
		return nil, nil, err
	}

	//Creating Refresh Token
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUuid
	rtClaims["user_id"] = user.Id
	rtClaims["username"] = user.Username
	rtClaims["name"] = user.Name
	rtClaims["surname"] = user.Surname
	rtClaims["patronymic"] = user.Patronymic
	rtClaims["exp"] = td.RtExpires

	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(refreshKey))
	if err != nil {
		return nil, nil, err
	}

	return td, user, nil

}

func (a *AuthUseCase) ParseAcsToken(ctx context.Context, accessToken string) (*models.AccessDetails, error) {
	token, err := jwt.ParseWithClaims(accessToken, &AuthClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(signingKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*AuthClaims); ok && token.Valid {
		return &models.AccessDetails{
			AccessUUID: claims.TokenUUID.String(),
			UserId:     claims.UserID,
			Name:       claims.Name,
			Surname:    claims.Surname,
			Patronymic: claims.Patronymic,
		}, nil
	}

	return nil, auth.ErrInvalidAccessToken

}

func (a *AuthUseCase) ParseRefToken(ctx context.Context, refreshToken string) (string, *models.User2, error) {
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(refreshKey), nil
	})

	user := new(models.User2)

	//if there is an error, the token must have expired
	if err != nil {
		return "", nil, err
	}

	//is token valid?
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return "", nil, errors.Errorf("%v", http.StatusUnprocessableEntity)
	}

	//Since token is valid, get the uuid:
	claims, ok := token.Claims.(jwt.MapClaims) //the token claims should conform to MapClaims
	if ok && token.Valid {
		refreshUuid, ok := claims["refresh_uuid"].(string) //convert to string
		if !ok {
			return "", nil, errors.Errorf("%v", http.StatusUnprocessableEntity)
		}
		user.Username, ok = claims["username"].(string)
		if !ok {
			return "", nil, errors.Errorf("%v", http.StatusUnprocessableEntity)
		}
		user.Name, ok = claims["name"].(string)
		if !ok {
			return "", nil, errors.Errorf("%v", http.StatusUnprocessableEntity)
		}
		user.Surname, ok = claims["surname"].(string)
		if !ok {
			return "", nil, errors.Errorf("%v", http.StatusUnprocessableEntity)
		}
		user.Patronymic, ok = claims["patronymic"].(string)
		if !ok {
			return "", nil, errors.Errorf("%v", http.StatusUnprocessableEntity)
		}
		user.Id = uint64(claims["user_id"].(float64))
		if err != nil {
			logrus.Info("5")
			return "", nil, errors.Errorf("%v: ok is %v", http.StatusUnprocessableEntity, ok)
		}

		return refreshUuid, user, nil
	}

	return "", nil, err
}

func (a *AuthUseCase) DeleteTokens(ctx context.Context, tokensUUID ...string) (uint64, error) {

	return a.stg.DeleteToken(ctx, tokensUUID)

}

func (a *AuthUseCase) CreateAuth(ctx context.Context, userid uint64, td *models.TokenDetails) error {

	return a.stg.CreateAuth(ctx, userid, td)

}

func (a *AuthUseCase) LogOut(ctx context.Context, givenUUID ...string) (uint64, error) {

	return a.stg.DeleteToken(ctx, givenUUID)

}

/*func (a *AuthUseCase) NewRefreshToken() (string, error) {
	b := make([]byte, 32)

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	_, err := r.Read(b)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}*/
