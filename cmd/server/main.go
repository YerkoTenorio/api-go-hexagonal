package main

import (
	"log"

	"github.com/YerkoTenorio/api-go-hexagonal/modules/task/application"
	"github.com/YerkoTenorio/api-go-hexagonal/modules/task/infrastructure"
	"github.com/YerkoTenorio/api-go-hexagonal/modules/task/presentation"
	"github.com/YerkoTenorio/api-go-hexagonal/shared/config"
	"github.com/YerkoTenorio/api-go-hexagonal/shared/database"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// Cargar configuración
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Error cargando configuración:", err)
	}

	// Conectar a base de datos con GORM
	gormDB, err := database.NewGormDB(cfg)
	if err != nil {
		log.Fatal("Error conectando a la base de datos:", err)
	}
	defer gormDB.Close()

	// Crear repositorio con GORM
	taskRepository := infrastructure.NewGormTaskRepository(gormDB.GetDB())

	// Crear servicio de aplicación
	taskService := application.NewTaskService(taskRepository)

	// Crear handler con Fiber
	taskHandler := presentation.NewFiberTaskHandler(taskService)

	// Crear aplicación Fiber
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New())

	// Rutas de salud
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":    "ok",
			"framework": "fiber",
			"orm":       "gorm",
		})
	})

	// Configurar rutas de tareas
	presentation.SetupTaskRoutesFiber(app, taskHandler)

	// Iniciar servidor
	log.Printf("Servidor Fiber + GORM iniciado en %s:%s", cfg.Server.Host, cfg.Server.Port)
	log.Fatal(app.Listen(cfg.Server.Host + ":" + cfg.Server.Port))
}
