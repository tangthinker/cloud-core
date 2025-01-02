package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/tangthinker/cloud-core/api/storage"
	"log"
	"os"
)

func main() {

	app := fiber.New(fiber.Config{
		BodyLimit: 800 * 1024 * 1024, // 800MB
	})

	loggingFile, err := os.OpenFile("/home/tangthinker/code/go-projects/cloud-core/requests.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}

	app.Use(logger.New(logger.Config{
		Format:     "${time} ${method} ${path} - ${ip} - ${status} - ${latency}\nn",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "Local",
		Output:     loggingFile,
	}))

	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.SendString("Hello, World!")
	})

	apiGroup := app.Group("/api/v1/storage/")
	apiGroup.Use(cors.New())
	api := storage.NewApi("/home/tangthinker/Downloads")

	apiGroup.Post("/ls", api.LS)
	apiGroup.Get("/get", api.Get)
	apiGroup.Get("get-thumbnail", api.GetThumbnail)
	apiGroup.Post("/stat", api.Stat)
	apiGroup.Get("/download", api.Download)
	apiGroup.Post("/upload", api.Upload)

	log.Fatal(app.Listen(":9999"))

}
