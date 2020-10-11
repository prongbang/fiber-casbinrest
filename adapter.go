package fibercasbinrest

import (
	"encoding/json"
	"fmt"
)

const (
	RoleKey       = "roles"
	RoleAnonymous = "anonymous"
)

type Adapter interface {
	GetRoleByToken(reqToken string) []string
}

type roleAdapter struct {
	Secret []byte
}

func (r *roleAdapter) GetRoleByToken(reqToken string) []string {
	t := GetValue(reqToken, RoleKey, r.Secret)
	var roles []interface{}
	if b, err := json.Marshal(t); err == nil {
		_ = json.Unmarshal(b, &roles)
	}
	s := make([]string, len(roles))
	for i, v := range roles {
		s[i] = fmt.Sprint(v)
	}
	return s
}

func NewRoleAdapter(secret string) Adapter {
	return &roleAdapter{
		Secret: []byte(secret),
	}
}
