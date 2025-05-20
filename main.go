package main

import (
	"flag"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/tangthinker/cloud-core/api/storage"
	"github.com/tangthinker/cloud-core/config"
	"github.com/tangthinker/cloud-core/internal/middleware"
	"github.com/tangthinker/user-center/pkg"
	"log"
	"os"
)

var configPath = flag.String("c", "config.toml", "config path")

func main() {

	flag.Parse()

	cnf := config.NewConfig(*configPath)
	logFilePath := cnf.LogFilePath()
	storagePath := cnf.StoragePath()
	rootPath := cnf.RootPath()
	fmt.Println("logFilePath:", logFilePath)
	fmt.Println("storagePath:", storagePath)
	fmt.Println("rootPath:", rootPath)

	app := fiber.New(fiber.Config{
		BodyLimit: 800 * 1024 * 1024, // 800MB
	})

	loggingFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}

	app.Use(logger.New(logger.Config{
		Format:     "${time} ${method} ${path} - ${ip} - ${status} - ${latency}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "Local",
		Output:     loggingFile,
	}))

	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.SendString("Hello, World!")
	})

	apiGroup := app.Group("/api/v1/storage/", middleware.TokenValid)

	api := storage.NewApi(storagePath)

	apiGroup.Post("/ls", api.LS)
	apiGroup.Get("/get", api.Get)
	apiGroup.Get("/get-thumbnail", api.GetThumbnail)
	apiGroup.Post("/stat", api.Stat)
	apiGroup.Get("/download", api.Download)
	apiGroup.Post("/upload", api.Upload)

	apiGroup.Get("/m3u8-state", api.Trans2M3U8)

	authGroup := app.Group("/api/v1/")
	pkg.RegisterUserCenter(authGroup, rootPath)

	log.Fatal(app.Listen(":" + cnf.ServicePort()))

}
