package services

import (
	"time"
	"webapp-go/webapp/config"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type BearerService interface {
	GenerateToken(id uuid.UUID) (string, error)
	ValidateToken(token string) (uuid.UUID, error)
	RefreshToken(token string) (string, error)
	RevokeToken(token string) error
}

type bearerService struct {
	cfg config.Config
}

func NewBearerService(cfg config.Config) BearerService {
	return bearerService{cfg}
}

func (this bearerService) GenerateToken(id uuid.UUID) (t string, err error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
    claims["id"] = id.String()
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	t, err = token.SignedString([]byte(this.cfg.JWT.Secret))

	return
}

func (this bearerService) ValidateToken(token string) (id uuid.UUID, err error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(this.cfg.JWT.Secret), nil
	})
	if err != nil {
		return
	}

	claims := t.Claims.(jwt.MapClaims)
    id, err = uuid.Parse(claims["id"].(string))

	return
}

func (this bearerService) RefreshToken(token string) (t string, err error) {
	p, err := this.ValidateToken(token)
	if err != nil {
		return
	}

	t, err = this.GenerateToken(p)

	return
}

func (this bearerService) RevokeToken(token string) (err error) {
	// This is a dummy implementation. In a real-world application, you would
	// need to store the revoked tokens in a database or cache.
	return
}
