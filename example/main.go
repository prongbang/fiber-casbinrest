package main

import (
	"log"

	"github.com/casbin/casbin/v2"
	mongodbadapter "github.com/casbin/mongodb-adapter/v3"
	"github.com/gofiber/fiber/v2"
	fibercasbinrest "github.com/prongbang/fiber-casbinrest"
)

func main() {
	a, _ := mongodbadapter.NewAdapter("127.0.0.1:27017")

	e, _ := casbin.NewEnforcer("example/model.conf", a)

	// Modify the policy.
	_, _ = e.AddPolicy("anonymous", "/login", "GET")
	_, _ = e.AddPolicy("admin", "/admin", "(GET)|(POST)")
	_, _ = e.AddPolicy("admin", "/admin/user/:id", "GET")
	_, _ = e.AddPolicy("admin", "/manage/*", "GET")
	//e.RemovePolicy(...)

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
