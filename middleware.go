package fibercasbinrest

import (
	"github.com/casbin/casbin/v2"
	"github.com/gofiber/fiber/v2"
	"log"
	"strings"
)

type (
	Skipper func(*fiber.Ctx) bool
	Config  struct {
		Skipper  Skipper
		Enforcer *casbin.Enforcer
		Adapter  Adapter
	}
)

var (
	DefaultConfig = Config{
		Skipper: DefaultSkipper,
	}
)

func DefaultSkipper(*fiber.Ctx) bool {
	return false
}

func New(ce *casbin.Enforcer, adt Adapter) fiber.Handler {
	c := DefaultConfig
	c.Enforcer = ce
	c.Adapter = adt
	return middlewareWithConfig(c)
}

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
		if config.Skipper(c) || config.CheckPermissions(c) {
			return c.Next()
		}
		return fiber.ErrForbidden
	}
}

// GetRole gets the roles name from the request.
func (a *Config) GetRole(c *fiber.Ctx) []string {
	token := c.Get(fiber.HeaderAuthorization)
	authorization := strings.Split(token, "Bearer")
	if len(authorization) == 2 {
		return a.Adapter.GetRoleByToken(strings.TrimSpace(authorization[1]))
	}
	return []string{RoleAnonymous}
}

// CheckPermissions checks the role/path/method combination from the request.
func (a *Config) CheckPermissions(c *fiber.Ctx) bool {
	roles := a.GetRole(c)
	allowed := false
	for _, role := range roles {
		result, err := a.Enforcer.Enforce(strings.ToLower(role), c.Path(), c.Method())
		if  result && err == nil {
			allowed = true
		} else {
			log.Println(err)
		}
		log.Println(strings.ToLower(role), c.Path(), c.Method(), " -> ", allowed)
	}
	return allowed
}
