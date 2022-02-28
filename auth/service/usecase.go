package service

import (
	"Test_derictory/auth"
	"Test_derictory/models"
	"context"
	"crypto/sha1"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofrs/uuid"
	"time"
)

const (
	salt       = "bcb545454yh5p5HG"
	tokenTTL   = 15 * time.Minute
	refreshTTL = 24 * 7 * time.Hour
	signingKey = "QHhpZGlF2DG3SD3F3G2SDF3H4vCg=="
	refreshKey = "ds7B989umHJ98opi;m2"
)

type AuthClaims struct {
	jwt.StandardClaims
	User      *models.User2 `json:"user"`
	TokenUUID uuid.UUID     `json:"access_uuid"`
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

func (a *AuthUseCase) SignIn(ctx context.Context, username, password string) (*models.TokenDetails, uint64, error) {
	passwordHash := GeneratePasswordHash(password)

	user, err := a.repo.GetUser(ctx, username, passwordHash)
	if err != nil {
		return nil, 0, err
	}

	//Generate Access Token
	td := &models.TokenDetails{}
	td.AtExpires = time.Now().Add(tokenTTL).Unix()
	td.AccessUuid, err = uuid.NewV4()
	if err != nil {
		return nil, 0, err
	}

	td.RtExpires = time.Now().Add(refreshTTL).Unix()
	td.RefreshUuid, err = uuid.NewV4()
	if err != nil {
		return nil, 0, err
	}

	claims := &AuthClaims{
		User:      user,
		TokenUUID: td.AccessUuid,
		//AccessUUID: td.AccessUuid,
		StandardClaims: jwt.StandardClaims{
			// Токен перестанет быть валидным через 15 минут с момента его генерации
			ExpiresAt: td.AtExpires,
			// Время генерации токена
			IssuedAt: time.Now().Unix(),
			Subject:  string(rune(user.Id)),
		},
	}

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	td.AccessToken, err = at.SignedString([]byte(signingKey))
	if err != nil {
		return nil, 0, err
	}

	//Creating Refresh Token
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUuid
	rtClaims["user"] = user
	rtClaims["exp"] = td.RtExpires

	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(refreshKey))
	if err != nil {
		return nil, 0, err
	}

	return td, user.Id, nil

}

func (a *AuthUseCase) ParseToken(ctx context.Context, accessToken string) (*models.AccessDetails, error) {
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
			UserId:     claims.User.Id,
		}, nil
	}

	return nil, auth.ErrInvalidAccessToken

}

func (a *AuthUseCase) CreateAuth(ctx context.Context, userid uint64, td *models.TokenDetails) error {

	return a.stg.CreateAuth(ctx, userid, td)

}

func (a *AuthUseCase) LogOut(ctx context.Context, givenUUID string) (int64, error) {

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
