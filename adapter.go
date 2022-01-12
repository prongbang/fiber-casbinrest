package fibercasbinrest

const (
	// RoleKey default
	RoleKey = "roles"
	// RoleAnonymous anonymous
	RoleAnonymous = "anonymous"
)

// Adapter interface for implements GetRoleByToken
type Adapter interface {
	GetRoleByToken(reqToken string) ([]string, error)
}
