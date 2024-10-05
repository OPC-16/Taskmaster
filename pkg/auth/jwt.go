package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
   UserID int64 `json:"user_id"`
   jwt.RegisteredClaims
}

// it implements the jwt.Claims interface
func (c Claims) Valid() error {
   if c.ExpiresAt == nil || c.ExpiresAt.Time.Before(time.Now()) {
      return errors.New("token is expired")
   }

   if c.UserID == 0 {
      return errors.New("must provide user id")
   }

   return nil
}

func GenerateToken(userID int64, secretKey string, expirationTime time.Duration) (string, error) {
   claims := &Claims{
      UserID: userID,
      RegisteredClaims: jwt.RegisteredClaims{
         ExpiresAt: jwt.NewNumericDate(time.Now().Add(expirationTime)),
         IssuedAt: jwt.NewNumericDate(time.Now()),
      },
   }

   token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
   return token.SignedString([]byte(secretKey))
}

func ValidateToken(tokenString string, secretKey string) (*Claims, error) {
   token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
      return []byte(secretKey), nil
   })

   if err != nil {
      return nil, err
   }

   if claims, ok := token.Claims.(*Claims); ok && token.Valid {
      return claims, nil
   }

   return nil, errors.New("invalid token")
}
