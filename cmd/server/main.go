package main

import (
    "log"
    "net/http"

    "github.com/YerkoTenorio/api-go-hexagonal/modules/task/application"
    "github.com/YerkoTenorio/api-go-hexagonal/modules/task/infrastructure"
    "github.com/YerkoTenorio/api-go-hexagonal/modules/task/presentation"
    "github.com/YerkoTenorio/api-go-hexagonal/shared/config"
    "github.com/YerkoTenorio/api-go-hexagonal/shared/database"
    "github.com/gin-gonic/gin"
)

func main() {
	// cargar la configuracion
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error cargando configuracion: %v", err)
	}

	// conectar base de datos
	db, err := database.NewSQLiteDB(cfg)
	if err != nil {
		log.Fatalf("Error conectando base de datos: %v", err)
	}
	defer db.Close()

	// Inyección de dependencias (arquitectura hexagonal)
	// Repositorio (Puerto de salida)
	taskRepo := infrastructure.NewSQLiteTaskRepository(db)
	// Servicio de aplicacion (Casos de uso)
	taskService := application.NewTaskService(taskRepo)
	// Handler (Puerto de entrada)
	taskHandler := presentation.NewTaskHandler(taskService)

	// configurar gin
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}
	// crear router
	router := gin.Default()
	// Middleware basico
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

    // Configurar rutas
    presentation.SetupRoutes(router, taskHandler)

    // Endpoints de salud
    router.GET("/health", func(c *gin.Context) {
        dbKind := "sqlite"
        isRemote := false
        if cfg.Database.URL != "" {
            dbKind = "libSQL"
            isRemote = true
        }

        c.JSON(http.StatusOK, gin.H{
            "status": "ok",
            "env":    cfg.App.Environment,
            "server": cfg.ServerAddress(),
            "database": gin.H{
                "type":   dbKind,
                "remote": isRemote,
            },
        })
    })

    router.GET("/health/db", func(c *gin.Context) {
        err := db.GetDB().Ping()

        var version string
        if err == nil {
            _ = db.GetDB().QueryRow("SELECT sqlite_version()").Scan(&version)
        }

        statusCode := http.StatusOK
        if err != nil {
            statusCode = http.StatusServiceUnavailable
        }

        provider := "sqlite"
        if cfg.Database.URL != "" {
            provider = "libSQL"
        }

        var errMsg string
        if err != nil {
            errMsg = err.Error()
        }

        c.JSON(statusCode, gin.H{
            "status":   map[bool]string{true: "up", false: "down"}[err == nil],
            "version":  version,
            "remote":   cfg.Database.URL != "",
            "provider": provider,
            "error":    errMsg,
        })
    })
	// iniciar el servidor
	log.Printf("Servidor iniciado en %s", cfg.ServerAddress())
	if err := router.Run(cfg.GetServerAddress()); err != nil {
		log.Fatalf("Error iniciando el servidor: %v", err)
	}

}
