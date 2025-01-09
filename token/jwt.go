package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	tokenExpiresSeconds        = 60 * 60 * 24
	refreshTokenExpiresSeconds = 60 * 60 * 25
)

type (
	Config struct {
		TokenExpiresSeconds   int64  `json:"expires_seconds"`
		RefreshExpiresSeconds int64  `json:"refresh_expires_seconds"`
		Key                   string `json:"key"`
		Issuer                string `json:"issuer"`
		Subject               string `json:"subject"`
	}

	JWT struct {
		*Config
	}

	Account struct {
		ID                    string `json:"user_id"`                  // id
		Role                  int    `json:"role"`                     // 角色
		RefreshTokenExpiresAt int64  `json:"refresh_token_expires_at"` // 该字段不需要赋值，取配置文件中的时间，程序中更改容易跟access_token的过期时间冲突
	}

	claims struct {
		*Account
		jwt.RegisteredClaims
	}
)

// GetExpiredTime implements Tokener.
func (j *JWT) GetExpiredTime(token string) (time.Time, error) {
	var c claims

	_, err := jwt.ParseWithClaims(token, &c, func(t *jwt.Token) (interface{}, error) {
		return j.GetKey(), nil
	})
	if err != nil {
		return time.Time{}, nil
	}

	nd, err := c.GetExpirationTime()
	if err != nil {
		return time.Time{}, err
	}
	return nd.Time, nil
}

func (j *JWT) ParseToken(token string) (*Account, error) {
	t, err := jwt.ParseWithClaims(token, &claims{}, func(t *jwt.Token) (interface{}, error) {
		return j.GetKey(), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			if value, ok := t.Claims.(*claims); ok {
				return value.Account, err
			}
		}
		return nil, err
	}

	if value, ok := t.Claims.(*claims); ok {
		return value.Account, err
	}

	return nil, fmt.Errorf("invalid token")
}

func (j *JWT) GenerateToken(account Account) (string, error) {
	account.RefreshTokenExpiresAt = time.Now().Unix() + j.RefreshExpiresSeconds
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims{
		Account:          &account,
		RegisteredClaims: j.newRegisteredClaims(),
	})

	return token.SignedString(j.GetKey())
}

func (j *JWT) GetKey() any {
	return []byte(j.Key)
}

func (j *JWT) newRegisteredClaims() jwt.RegisteredClaims {
	return jwt.RegisteredClaims{
		Issuer:    j.Issuer,
		Subject:   j.Subject,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(j.TokenExpiresSeconds) * time.Second)),
		NotBefore: jwt.NewNumericDate(time.Now()),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}
}

func NewJWT(cfg *Config) (*JWT, error) {
	if cfg == nil {
		return nil, fmt.Errorf("empty jwt cnf")
	}

	if cfg.Key == "" {
		return nil, fmt.Errorf("empty key")
	}

	if cfg.TokenExpiresSeconds <= 0 {
		cfg.TokenExpiresSeconds = tokenExpiresSeconds
	}

	if cfg.RefreshExpiresSeconds <= 0 {
		cfg.RefreshExpiresSeconds = refreshTokenExpiresSeconds
	}

	if cfg.RefreshExpiresSeconds <= cfg.TokenExpiresSeconds {
		return nil, fmt.Errorf("refresh expires  must be bigger than token expires")
	}

	return &JWT{
		Config: cfg,
	}, nil
}
