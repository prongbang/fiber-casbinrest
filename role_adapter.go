package fibercasbinrest

type roleAdapter struct {
	Secret []byte
}

func (r *roleAdapter) GetRoleByToken(reqToken string) ([]string, error) {
	t, err := GetValue(reqToken, RoleKey, r.Secret)
	return ParseRoles(t), err
}

// NewRoleAdapter create adapter
func NewRoleAdapter(secret string) Adapter {
	return &roleAdapter{
		Secret: []byte(secret),
	}
}
