package main

import (
	"github.com/casbin/casbin/v2"
	mongodbadapter "github.com/casbin/mongodb-adapter/v3"
	"github.com/gofiber/fiber/v2"
	fibercasbinrest "github.com/prongbang/fiber-casbinrest"
	"log"
)

func main() {

	// Initialize a MongoDB adapter and use it in a Casbin enforcer:
	// The adapter will use the database named "casbin".
	// If it doesn't exist, the adapter will create it automatically.
	a, err := mongodbadapter.NewAdapter("127.0.0.1:27017")
	if err != nil {
		panic(err)
	}
	// Or you can use an existing DB "abc" like this:
	// The adapter will use the table named "casbin_rule".
	// If it doesn't exist, the adapter will create it automatically.
	// a := mongodbadapter.NewAdapter("127.0.0.1:27017/abc")

	e, err := casbin.NewEnforcer("example/model.conf", a)
	if err != nil {
		panic(err)
	}

	// Initialize a Redis adapter and use it in a Casbin enforcer:
	//a := redisadapter.NewAdapterWithPassword("tcp", "127.0.0.1:6379", "640c8dd3-54a7-42f5-9249-7d62dab232d2")
	//e, err := casbin.NewEnforcer("example/rbac_model.conf", a)
	//if err != nil {
	//	panic(err)
	//}

	// Check the permission.
	//r, err := e.Enforce("admin", "/admin", "POST")
	//log.Println(r)
	//log.Println(err)

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
