package presentation

import (
	"strconv"

	"github.com/YerkoTenorio/api-go-hexagonal/modules/task/application"
	"github.com/gofiber/fiber/v2"
)

// FiberTaskHandler maneja las peticiones HTTP con Fiber
type FiberTaskHandler struct {
	taskService application.TaskServiceInterface
}

// NewFiberTaskHandler crea una nueva instancia del handler de tareas con Fiber
func NewFiberTaskHandler(taskService application.TaskServiceInterface) *FiberTaskHandler {
	return &FiberTaskHandler{
		taskService: taskService,
	}
}

// CreateTaskRequest representa la estructura de la petición para crear una tarea
type FiberCreateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// UpdateTaskRequest representa la estructura de la petición para actualizar una tarea
type FiberUpdateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   *bool  `json:"completed,omitempty"`
}

// CreateTask maneja la creación de una nueva tarea con Fiber
func (h *FiberTaskHandler) CreateTask(c *fiber.Ctx) error {
	var req FiberCreateTaskRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request",
			"message": err.Error(),
		})
	}

	// Validación manual (Fiber no tiene validación automática como Gin)
	if req.Title == "" || req.Description == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Validation failed",
			"message": "title and description are required",
		})
	}

	task, err := h.taskService.CreateTask(c.Context(), req.Title, req.Description)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Error creating task",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Task created successfully",
		"data":    task,
	})
}

// GetAllTasks obtiene todas las tareas con Fiber
func (h *FiberTaskHandler) GetAllTasks(c *fiber.Ctx) error {
	tasks, err := h.taskService.GetAllTasks(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Error getting tasks",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Tasks retrieved successfully",
		"data":    tasks,
		"count":   len(tasks),
	})
}

// GetTask obtiene una tarea por su ID con Fiber
func (h *FiberTaskHandler) GetTask(c *fiber.Ctx) error {
	idStr := c.Params("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid ID",
			"message": "ID must be a positive integer",
		})
	}

	task, err := h.taskService.GetTaskByID(c.Context(), int(id))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Error getting task",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Task retrieved successfully",
		"data":    task,
	})
}

// UpdateTask actualiza una tarea existente con Fiber
func (h *FiberTaskHandler) UpdateTask(c *fiber.Ctx) error {
	idStr := c.Params("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid ID",
			"message": "ID must be a positive integer",
		})
	}

	var req FiberUpdateTaskRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request",
			"message": err.Error(),
		})
	}

	task, err := h.taskService.UpdateTask(c.Context(), int(id), req.Title, req.Description, req.Completed)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Error updating task",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Task updated successfully",
		"data":    task,
	})
}

// DeleteTask elimina una tarea existente con Fiber
func (h *FiberTaskHandler) DeleteTask(c *fiber.Ctx) error {
	idStr := c.Params("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid ID",
			"message": "ID must be a positive integer",
		})
	}

	err = h.taskService.DeleteTask(c.Context(), int(id))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Error deleting task",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Task deleted successfully",
	})
}

// GetTaskByStatus obtiene todas las tareas filtradas por estado con Fiber
func (h *FiberTaskHandler) GetTaskByStatus(c *fiber.Ctx) error {
	completedStr := c.Query("completed")
	completed, err := strconv.ParseBool(completedStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid completed parameter",
			"message": "completed must be a boolean value",
		})
	}

	tasks, err := h.taskService.GetTasksByStatus(c.Context(), completed)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Error getting tasks",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Tasks retrieved successfully",
		"data":    tasks,
	})
}
