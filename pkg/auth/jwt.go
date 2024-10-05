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
