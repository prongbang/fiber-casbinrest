package fibercasbinrest

import (
	"log"

	"github.com/dgrijalva/jwt-go"
)

// GetValue for get payload from JWT
func GetValue(reqToken string, key string, secretKey []byte) interface{} {
	token, err := jwt.Parse(reqToken, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		log.Println(err)
		return ""
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		uid := claims[key]
		return uid
	}
	log.Println(claims.Valid().Error())
	return ""
}
