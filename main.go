package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tangthinker/cloud-core/api/storage"
	"log"
)

func main() {

	app := fiber.New()

	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.SendString("Hello, World!")
	})

	apiGroup := app.Group("/api/v1/storage/")
	api := storage.NewApi("/Users/tal/code/GoProject/cloud-core")

	apiGroup.Post("/ls", api.LS)
	apiGroup.Post("/get", api.Get)
	apiGroup.Post("/stat", api.Stat)
	apiGroup.Get("/download", api.Download)
	apiGroup.Post("/upload", api.Upload)

	log.Fatal(app.Listen(":9999"))

}
