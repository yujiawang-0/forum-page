package auth

// referenced from: https://medium.com/@cheickzida/golang-implementing-jwt-token-authentication-bba9bfd84d60

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
)
var secretKey = []byte("secret-key")

type Claims struct {
	UserID	int		`json:"user_id"`
	Role 	string	`json:"role"`
	jwt.RegisteredClaims
}

// create a new JWT token using jwt.NewWithClaims()
// signing method = HS256
// sign token with a secret key and return gnerated token as string
func createToken(userID int, role string) (string, error) {
	claims := Claims{
		UserID: userID,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},

	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

   tokenString, err := token.SignedString(secretKey)
   if err != nil {
   return "", err
   }

 return tokenString, nil
}

// verify JWT toekn authenticity 
// we use the jwt.Parse() function to parse and verify the token. We provide a callback function to retrieve the secret key used for signing the token. If the token is valid, we continue processing the request; otherwise, we return an error indicating that the token is invalid.

func verifyToken(tokenString string) (*Claims, error) {
   token, err := jwt.ParseWithClaims(
	tokenString, 
	&Claims{},
	func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
      return secretKey, nil
   })
  
   if err != nil {
      return nil, err
   }
   
   claims, ok := token.Claims.(*Claims)
   if !ok || !token.Valid {
      return nil, errors.New("invalid token")
   }
  
   return claims, nil
}