package presentation

import (
	"net/http"
	"strconv"

	"github.com/YerkoTenorio/api-go-hexagonal/modules/task/application"
	"github.com/gin-gonic/gin"
)

// TaskHandler maneja las peticiones HTTP relacionadas con tareas
type TaskHandler struct {
	taskService application.TaskServiceInterface
}

// NewTaskHandler crea una nueva instancia del handler de tareas
func NewTaskHandler(taskService application.TaskServiceInterface) *TaskHandler {
	return &TaskHandler{
		taskService: taskService,
	}
}

// CreateTaskRequest representa la estructura de la peticion para crear una tarea
type CreateTaskRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
}

// UpdateTaskRequest representa la estructura de la peticion para actualizar una tarea

type UpdateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed,omitempty"`
}

// CreateTask maneja la creacion de una nueva tarea
// @Summary Crea una nueva tarea
// @Description Crea una nueva tarea con el titulo y descripcion proporcionados
// @Tags tareas
// @Accept json
// @Produce json
// @Param task body CreateTaskRequest true "Dato de la tarea a crear"
// @Success 201 {object} entities.Task
// @Failure 400 {object} gin.H
// @Router /tasks [post]

func (h *TaskHandler) CreateTask(c *gin.Context) {
	var req CreateTaskRequest
	// Gin automaticamente valida y bindea el JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": err.Error(),
		})
		return
	}

	// Crear la tarea usando el servicio
	task, err := h.taskService.CreateTask(c.Request.Context(), req.Title, req.Description)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Error creating task",
			"message": err.Error(),
		})
		return
	}
	// Respuesta exitosa
	c.JSON(http.StatusCreated, gin.H{
		"message": "Task created successfully",
		"data":    task,
	})
}

// GetAllTasks obtiene todas las tareas
// @Summary Obtiene todas las tareas
// @Description Obtiene todas las tareas almacenadas en el sistema
// @Tags tareas
// @Produce json
// @Success 200 {object} []entities.Task
// @Failure 500 {object} gin.H
// @Router /tasks [get]
func (h *TaskHandler) GetAllTasks(c *gin.Context) {
	// Obtener todas las tareas usando el servicio
	tasks, err := h.taskService.GetAllTasks(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error getting tasks",
			"message": err.Error(),
		})
		return
	}

	// Respuesta exitosa
	c.JSON(http.StatusOK, gin.H{
		"message": "Tasks retrieved successfully",
		"data":    tasks,
		"count":   len(tasks),
	})
}

// GetTask obtiene una tarea por su ID
// @Summary Obtiene una tarea por ID
// @Description Obtiene los detalles de una tarea especifica por su ID
// @Tags tareas
// @Produce json
// @Param id path int true "ID de la tarea"
// @Success 200 {object} entities.Task
// @Failure 400 {object} gin.H
// @Failure 404 {object} gin.H
// @Router /tasks/{id} [get]
func (h *TaskHandler) GetTask(c *gin.Context) {
	// Gin facilita obtener parametros de la URL
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid ID",
			"message": "ID must be a positive integer",
		})
		return
	}
	// Obtener la tarea usando el servicio
	task, err := h.taskService.GetTaskByID(c.Request.Context(), int(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Error getting task",
			"message": err.Error(),
		})
		return
	}
	// Respuesta exitosa
	c.JSON(http.StatusOK, gin.H{
		"message": "Task retrieved successfully",
		"data":    task,
	})

}

// UpdateTask actualiza una tarea existente
// @Summary Actualiza una tarea existente
// @Description Actualiza los detalles de una tarea especifica por su ID
// @Tags tareas
// @Accept json
// @Produce json
// @Param id path int true "ID de la tarea"
// @Param task body UpdateTaskRequest true "Dato de la tarea a actualizar"
// @Success 200 {object} entities.Task
// @Failure 400 {object} gin.H
// @Failure 404 {object} gin.H
// @Router /tasks/{id} [put]

func (h *TaskHandler) UpdateTask(c *gin.Context) {
	// Obtener el ID de la URL
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid ID",
			"message": "ID must be a positive integer",
		})
		return
	}

	var req UpdateTaskRequest
	// Bindear el JSON (sin validacion required porque son campos opcionales)
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": err.Error(),
		})
		return
	}
	// Actualizar la tarea usando el servicio
	task, err := h.taskService.UpdateTask(c.Request.Context(), int(id), req.Title, req.Description, &req.Completed)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Error updating task",
			"message": err.Error(),
		})
		return
	}
	// Respuesta exitosa
	c.JSON(http.StatusOK, gin.H{
		"message": "Task updated successfully",
		"data":    task,
	})
}

// DeleteTask elimina una tarea existente
// @Summary Elimina una tarea existente
// @Description Elimina una tarea especifica por su ID
// @Tags tareas
// @Produce json
// @Param id path int true "ID de la tarea"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 404 {object} gin.H
// @Router /tasks/{id} [delete]
func (h *TaskHandler) DeleteTask(c *gin.Context) {
	// Obtener el ID de la URL
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid ID",
			"message": "ID must be a positive integer",
		})
		return
	}
	// Eliminar la tarea usando el servicio
	err = h.taskService.DeleteTask(c.Request.Context(), int(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Error deleting task",
			"message": err.Error(),
		})
		return
	}
	// Respuesta exitosa
	c.JSON(http.StatusOK, gin.H{
		"message": "Task deleted successfully",
	})
}

// GetTaskByStatus obtiene todas las tareas filtradas por estado
// @Summary Obtiene todas las tareas filtradas por estado
// @Description Obtiene todas las tareas que esten completadas o no segun el estado
// @Tags tareas
// @Produce json
// @Param completed query boolean false "Filtrar por estado completada"
// @Success 200 {object} []entities.Task
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /tasks/status [get]

func (h *TaskHandler) GetTaskByStatus(c *gin.Context) {
	// Obtener el parametro de consulta completed
	completedStr := c.Query("completed")
	completed, err := strconv.ParseBool(completedStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid completed parameter",
			"message": "completed must be a boolean value",
		})
		return
	}
	// Obtener todas las tareas usando el servicio
	tasks, err := h.taskService.GetTasksByStatus(c.Request.Context(), completed)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error getting tasks",
			"message": err.Error(),
		})
		return
	}
	// Respuesta exitosa
	c.JSON(http.StatusOK, gin.H{
		"message": "Tasks retrieved successfully",
		"data":    tasks,
	})
}
