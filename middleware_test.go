package fibercasbinrest_test

import (
	"fmt"
	"github.com/casbin/casbin/v2"
	"github.com/gofiber/fiber/v2"
	fibercasbinrest "github.com/prongbang/fiber-casbinrest"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

type redisAdapter struct {
}

func NewRedisAdapter() fibercasbinrest.Adapter {
	return &redisAdapter{}
}

const mockAdminToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"

func (r *redisAdapter) GetRoleByToken(reqToken string) ([]string, error) {
	role := "anonymous"
	if reqToken == mockAdminToken {
		role = "admin"
	} else if reqToken == "TOKEN_DBA" {
		role = "dba"
	}
	return []string{role}, nil
}

var adapter fibercasbinrest.Adapter

func init() {
	adapter = NewRedisAdapter()
}

func TestRoleAdminStatusOK(t *testing.T) {
	// Given
	ce, _ := casbin.NewEnforcer("example/auth_model.conf", "example/policy.csv")
	e := fiber.New()
	e.Use(fibercasbinrest.New(ce, adapter))
	e.Get("/", func(c *fiber.Ctx) error {
		return c.Status(http.StatusOK).JSON("OK")
	})
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", mockAdminToken))

	// When
	res, _ := e.Test(req, 10000)

	// Then
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestRoleDbaStatusOK(t *testing.T) {
	// Given
	ce, _ := casbin.NewEnforcer("example/auth_model.conf", "example/policy.csv")
	e := fiber.New()
	e.Use(fibercasbinrest.New(ce, adapter))
	e.Post("/dba", func(c *fiber.Ctx) error {
		return c.Status(http.StatusOK).JSON("OK")
	})
	req := httptest.NewRequest(http.MethodPost, "/dba", nil)
	req.Header.Set("Authorization", "Bearer TOKEN_DBA")

	// When
	res, _ := e.Test(req, 10000)

	// Then
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestRoleAdminStatusForbidden(t *testing.T) {
	// Given
	ce, _ := casbin.NewEnforcer("example/auth_model.conf", "example/policy.csv")
	e := fiber.New()
	e.Use(fibercasbinrest.New(ce, adapter))
	e.Get("/", func(c *fiber.Ctx) error {
		return c.Status(http.StatusOK).JSON("OK")
	})
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", "Mock Token"))

	// When
	res, _ := e.Test(req, 10000)

	// Then
	assert.Equal(t, http.StatusForbidden, res.StatusCode)
}

func TestRoleAnonymousWithoutTokenStatusForbidden(t *testing.T) {
	// Given
	ce, _ := casbin.NewEnforcer("example/auth_model.conf", "example/policy.csv")
	e := fiber.New()
	e.Use(fibercasbinrest.New(ce, adapter))
	e.Get("/logout", func(c *fiber.Ctx) error {
		return c.Status(http.StatusOK).JSON("OK")
	})
	req := httptest.NewRequest(http.MethodGet, "/logout", nil)

	// When
	res, _ := e.Test(req, 10000)

	// Then
	assert.Equal(t, http.StatusForbidden, res.StatusCode)
}

func TestRoleAnonymousTokenStatusOK(t *testing.T) {
	// Given
	ce, _ := casbin.NewEnforcer("example/auth_model.conf", "example/policy.csv")
	e := fiber.New()
	e.Use(fibercasbinrest.New(ce, adapter))
	e.Get("/login", func(c *fiber.Ctx) error {
		return c.Status(http.StatusOK).JSON("OK")
	})
	req := httptest.NewRequest(http.MethodGet, "/login", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", "Mock Token"))

	// When
	res, _ := e.Test(req, 10000)

	// Then
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestRoleAdminByJWTStatusOK(t *testing.T) {
	// Given
	secret := "test"
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJyb2xlcyI6WyJBRE1JTiIsIlVTRVIiXX0.7fOkJTmCRWJTVR5lLh0He_wtUTEUfWEFkrcfoArgPIw"
	ce, _ := casbin.NewEnforcer("example/auth_model.conf", "example/policy.csv")
	e := fiber.New()
	e.Use(fibercasbinrest.NewDefault(ce, secret))
	e.Get("/", func(c *fiber.Ctx) error {
		return c.Status(http.StatusOK).JSON("OK")
	})
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	// When
	res, _ := e.Test(req, 10000)

	// Then
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestRoleAdminByJWTTokenExpired(t *testing.T) {
	// Given
	secret := "secret"
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NDEsInJvbGVzIjpbIkFETUlOIiwiVVNFUiJdfQ.P7B4nnVuw6FUscVtKLUn011Q0iZssO7LEr_o7d8nprE"
	ce, _ := casbin.NewEnforcer("example/auth_model.conf", "example/policy.csv")
	e := fiber.New()
	e.Use(fibercasbinrest.NewDefault(ce, secret))
	e.Get("/", func(c *fiber.Ctx) error {
		return c.Status(http.StatusOK).JSON("OK")
	})
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	// When
	res, _ := e.Test(req, 10000)

	// Then
	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
}
