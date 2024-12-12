package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/tangthinker/cloud-core/api/storage"
	"log"
)

func main() {

	app := fiber.New()

	app.Static("/", "./public")

	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.SendString("Hello, World!")
	})

	apiGroup := app.Group("/api/v1/storage/")
	apiGroup.Use(cors.New())
	api := storage.NewApi("/Users/tal/Downloads")

	apiGroup.Post("/ls", api.LS)
	apiGroup.Get("/get", api.Get)
	apiGroup.Get("get-thumbnail", api.GetThumbnail)
	apiGroup.Post("/stat", api.Stat)
	apiGroup.Get("/download", api.Download)
	apiGroup.Post("/upload", api.Upload)

	log.Fatal(app.Listen(":9999"))

}
