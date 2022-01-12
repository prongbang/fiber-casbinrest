package fibercasbinrest

import (
	"net/http"
	"strings"

	"github.com/casbin/casbin/v2"
	"github.com/gofiber/fiber/v2"
)

type (
	// Skipper middleware function
	Skipper func(*fiber.Ctx) bool
	// Config middleware model
	Config struct {
		Skipper  Skipper
		Enforcer *casbin.Enforcer
		Adapter  Adapter
	}
)

var (
	// DefaultConfig initial default config
	DefaultConfig = Config{
		Skipper: DefaultSkipper,
	}
)

// DefaultSkipper create default skipper
func DefaultSkipper(*fiber.Ctx) bool {
	return false
}

// New create middleware
func New(ce *casbin.Enforcer, adt Adapter) fiber.Handler {
	c := DefaultConfig
	c.Enforcer = ce
	c.Adapter = adt
	return middlewareWithConfig(c)
}

// NewDefault create middleware
func NewDefault(ce *casbin.Enforcer, secret string) fiber.Handler {
	c := DefaultConfig
	c.Enforcer = ce
	c.Adapter = NewRoleAdapter(secret)
	return middlewareWithConfig(c)
}

func middlewareWithConfig(config Config) fiber.Handler {
	if config.Skipper == nil {
		config.Skipper = DefaultConfig.Skipper
	}
	return func(c *fiber.Ctx) error {
		pass, err := config.CheckPermissions(c)
		if config.Skipper(c) || (pass && err == nil) {
			return c.Next()
		}
		if err != nil && strings.ToLower(err.Error()) == "token is expired" {
			return c.Status(http.StatusUnauthorized).
				JSON(fiber.Map{"message": err.Error()})
		}
		return c.Status(http.StatusForbidden).
			JSON(fiber.Map{"message": http.StatusText(http.StatusForbidden)})
	}
}

// GetRole gets the roles name from the request.
func (a *Config) GetRole(c *fiber.Ctx) ([]string, error) {
	token := c.Get(fiber.HeaderAuthorization)
	authorization := strings.Split(token, "Bearer")
	if len(authorization) == 2 {
		return a.Adapter.GetRoleByToken(strings.TrimSpace(authorization[1]))
	}
	return []string{RoleAnonymous}, nil
}

// CheckPermissions checks the role/path/method combination from the request.
func (a *Config) CheckPermissions(c *fiber.Ctx) (bool, error) {
	roles, err := a.GetRole(c)
	allowed := false
	for _, role := range roles {
		result, e := a.Enforcer.Enforce(strings.ToLower(role), c.Path(), c.Method())
		if result && e == nil {
			allowed = true
		}
	}
	return allowed, err
}
