package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Manager interface {
	Generate(userID string) (string, error)
	Verify(tokenStr string) (*jwt.RegisteredClaims, error)
}

type jwtManager struct {
	secretKey     string
	tokenDuration time.Duration
}

func NewJWTManager(secretKey string, tokenDuration time.Duration) Manager {
	return &jwtManager{secretKey, tokenDuration}
}

// Generates a JSON Web Token for the given user ID.
// Returns the token string or an error if generation fails.
func (j *jwtManager) Generate(userID string) (string, error) {
	claims := &jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.tokenDuration)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

// Verifies the given token string and returns the claims if valid.
// Returns an error if the token is invalid or expired.
func (j *jwtManager) Verify(tokenStr string) (*jwt.RegisteredClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&jwt.RegisteredClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(j.secretKey), nil
		},
	)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return nil, jwt.ErrInvalidKeyType
	}
	return claims, nil
}
