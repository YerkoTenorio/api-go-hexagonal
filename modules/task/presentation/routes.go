package presentation

import (
	"github.com/gin-gonic/gin"
)

// SetupTaskRoutes configura todas las rutas relacionadas con tareas
func SetupTaskRoutes(router *gin.Engine, taskHandler *TaskHandler) {
	// Grupo de rutas para tareas
	taskGroup := router.Group("/api/v1/tasks")
	{
		// POST /api/v1/tasks - Crear nueva tarea
		taskGroup.POST("", taskHandler.CreateTask)

		// GET /api/v1/tasks - Obtener todas las tareas
		taskGroup.GET("", taskHandler.GetAllTasks)

		// GET /api/v1/tasks/:id - Obtener tarea por ID
		taskGroup.GET("/:id", taskHandler.GetTask)

		// PUT /api/v1/tasks/:id - Actualizar tarea
		taskGroup.PUT("/:id", taskHandler.UpdateTask)

		// DELETE /api/v1/tasks/:id - Eliminar tarea
		taskGroup.DELETE("/:id", taskHandler.DeleteTask)

		// GET /api/v1/tasks/status/:status - Obtener tareas por estado
		taskGroup.GET("/status/:status", taskHandler.GetTaskByStatus)
	}
}

// SetupRoutes configura todas las rutas de la aplicación
func SetupRoutes(router *gin.Engine, taskHandler *TaskHandler) {
    // Configurar rutas de tareas
    SetupTaskRoutes(router, taskHandler)
}
