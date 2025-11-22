package presentation

import (
	"github.com/gofiber/fiber/v2"
)

// SetupTaskRoutesFiber configura las rutas de tareas para Fiber
func SetupTaskRoutesFiber(app *fiber.App, handler *FiberTaskHandler) {
	// Grupo de rutas para tareas
	tasks := app.Group("/tasks")

	// CRUD básico
	tasks.Post("/", handler.CreateTask)
	tasks.Get("/", handler.GetAllTasks)
	tasks.Get("/:id", handler.GetTask)
	tasks.Put("/:id", handler.UpdateTask)
	tasks.Delete("/:id", handler.DeleteTask)

	// Rutas adicionales
	tasks.Get("/status", handler.GetTaskByStatus)
}
