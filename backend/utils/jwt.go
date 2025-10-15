package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
	jwt "github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var jwtSecret = []byte("firma_secretisima_del_token")

type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateToken(userID primitive.ObjectID, email string, role string) (string, error) {
	claims := Claims{
		UserID: userID.Hex(),
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) { //por cada intento de acceso a una funcion privada de la app se llama a este metodo
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok { //con esta funcion nos aseguramos que el contenido del token este seguro, ya que se verifica que el mismo este firmado con algun metodo de HMAC y no haya sido alterado el header o el payload por algun usuario
			return nil, errors.New("")
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
