# Casbin RESTful Adapter on Fiber Web Framework

Casbin RESTful adapter for Casbin [https://github.com/casbin/casbin](https://github.com/casbin/casbin)

[![Build Status](http://img.shields.io/travis/prongbang/fiber-casbinrest.svg)](https://travis-ci.org/prongbang/fiber-casbinrest)
[![Codecov](https://img.shields.io/codecov/c/github/prongbang/fiber-casbinrest.svg)](https://codecov.io/gh/prongbang/fiber-casbinrest)
[![Go Report Card](https://goreportcard.com/badge/github.com/prongbang/fiber-casbinrest)](https://goreportcard.com/report/github.com/prongbang/fiber-casbinrest)

## Install

```
go get github.com/prongbang/fiber-casbinrest
```

## Config

- model.conf

```editorconfig
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = r.sub == p.sub && (keyMatch(r.obj, p.obj) || keyMatch2(r.obj, p.obj)) && (r.act == p.act || regexMatch(r.act, p.act))
```

- support policy

```
p, admin, /user/*, (GET)|(POST)
p, anonymous, /login, (GET)
p, admin, /admin/user/:id, (GET)|(POST)
```

## JWT

```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyLCJyb2xlcyI6WyJBRE1JTiJdfQ.oW8uC8uyL4nZSjcDGRkW3ZHoEoHShPD7ft0cppgvQe4
```

- payload

```json
{
  "sub": "1234567890",
  "name": "John Doe",
  "iat": 1516239022,
  "roles": [
    "ADMIN"
  ]
}
```

## Usage with File and JWT

```go
import (
	"github.com/casbin/casbin/v2"
	"github.com/gofiber/fiber/v2"
	fibercasbinrest "github.com/prongbang/fiber-casbinrest"
	"log"
)

func main() {
    e, _ := casbin.NewEnforcer("example/auth_model.conf", "example/policy.csv")
    
    app := fiber.New()
    app.Use(fibercasbinrest.NewDefault(e, "secret"))
    
    app.Post("/admin", func(c *fiber.Ctx) error {
        return c.SendString("Hello, World 👋!")
    })
    
    app.Get("/admin/user/:id", func(c *fiber.Ctx) error {
        return c.SendString("Hello, User 👋!")
    })
    
    app.Get("/manage/:id", func(c *fiber.Ctx) error {
        return c.SendString("Hello, Manage 👋!")
    })
    
    app.Get("/login", func(c *fiber.Ctx) error {
        return c.SendString("Hello, login 👋!")
    })
    
    log.Fatal(app.Listen(":3000"))
}
```

## Usage with Custom adapter and JWT

```go
import (
	"github.com/casbin/casbin/v2"
	"github.com/gofiber/fiber/v2"
	fibercasbinrest "github.com/prongbang/fiber-casbinrest"
	"log"
)

type redisAdapter struct {
}

func NewRedisAdapter() fibercasbinrest.Adapter {
	return &redisAdapter{}
}

const mockAdminToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"

func (r *redisAdapter) GetRoleByToken(reqToken string) []string {
    // Validate example not use on production
	role := "anonymous"
	if reqToken == mockAdminToken {
		role = "admin"
	} else if reqToken == "TOKEN_DBA" {
		role = "dba"
	}
	return []string{role}
}

func main() {
    adapter := NewRedisAdapter()
    e, _ := casbin.NewEnforcer("example/auth_model.conf", "example/policy.csv")
    
    app := fiber.New()
    app.Use(fibercasbinrest.New(e, adapter))
    
    app.Post("/admin", func(c *fiber.Ctx) error {
        return c.SendString("Hello, World 👋!")
    })
    
    app.Get("/admin/user/:id", func(c *fiber.Ctx) error {
        return c.SendString("Hello, User 👋!")
    })
    
    app.Get("/manage/:id", func(c *fiber.Ctx) error {
        return c.SendString("Hello, Manage 👋!")
    })
    
    app.Get("/login", func(c *fiber.Ctx) error {
        return c.SendString("Hello, login 👋!")
    })
    
    log.Fatal(app.Listen(":3000"))
}
```

## Usage with MongoDB

```go
import (
	"github.com/casbin/casbin/v2"
	mongodbadapter "github.com/casbin/mongodb-adapter/v3"
	"github.com/gofiber/fiber/v2"
	fibercasbinrest "github.com/prongbang/fiber-casbinrest"
	"log"
)

func main() {
    a, _ := mongodbadapter.NewAdapter("127.0.0.1:27017")
    e, _ := casbin.NewEnforcer("model.conf", a)
    
    // Modify the policy.
    _, _ = e.AddPolicy("anonymous", "/login", "GET")
    _, _ = e.AddPolicy("admin", "/admin", "(GET)|(POST)")
    _, _ = e.AddPolicy("admin", "/admin/user/:id", "GET")
    _, _ = e.AddPolicy("admin", "/manage/*", "GET")
    
    // Save the policy back to DB.
    _ = e.SavePolicy()
    
    // Load the policy from DB.
    _ = e.LoadPolicy()
    
    app := fiber.New()
    app.Use(fibercasbinrest.NewDefault(e, "secret"))
    
    app.Post("/admin", func(c *fiber.Ctx) error {
        return c.SendString("Hello, World 👋!")
    })
    
    app.Get("/admin/user/:id", func(c *fiber.Ctx) error {
        return c.SendString("Hello, User 👋!")
    })
    
    app.Get("/manage/:id", func(c *fiber.Ctx) error {
        return c.SendString("Hello, Manage 👋!")
    })
    
    app.Get("/login", func(c *fiber.Ctx) error {
        return c.SendString("Hello, login 👋!")
    })
    
    log.Fatal(app.Listen(":3000"))
}
```

### Functions in matchers

- https://casbin.org/docs/en/function#how-to-add-a-customized-function

### Editor

- https://casbin.org/editor/