package fibercasbinrest

import (
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt"
)

// Verify JWT
func Verify(token string, secret []byte) bool {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		return false
	}
	return t.Valid
}

// ParseToken for validate JWT
func ParseToken(token string, secret []byte) (*jwt.Token, error) {
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
}

// GetValue for get payload from JWT
func GetValue(reqToken string, key string, secretKey []byte) (interface{}, error) {
	token, err := ParseToken(reqToken, secretKey)
	if err != nil {
		return "", err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		return claims[key], nil
	}
	return "", claims.Valid()
}

// ParseRoles interface to string array
func ParseRoles(roles interface{}) []string {
	var rs []interface{}
	if b, err := json.Marshal(roles); err == nil {
		_ = json.Unmarshal(b, &rs)
	}
	s := make([]string, len(rs))
	for i, v := range rs {
		s[i] = fmt.Sprint(v)
	}
	return s
}
