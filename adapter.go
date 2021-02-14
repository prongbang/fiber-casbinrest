package fibercasbinrest

const (
	// RoleKey default
	RoleKey = "roles"
	// RoleAnonymous anonymous
	RoleAnonymous = "anonymous"
)

// Adapter interface for implements GetRoleByToken
type Adapter interface {
	GetRoleByToken(reqToken string) []string
}

type roleAdapter struct {
	Secret []byte
}

func (r *roleAdapter) GetRoleByToken(reqToken string) []string {
	t := GetValue(reqToken, RoleKey, r.Secret)
	return ParseRoles(t)
}

// NewRoleAdapter create adapter
func NewRoleAdapter(secret string) Adapter {
	return &roleAdapter{
		Secret: []byte(secret),
	}
}
