package application

import (
	"context"

	"github.com/YerkoTenorio/api-go-hexagonal/modules/task/domain"
)

//go:generate mockgen -source=interfaces.go -destination=mocks/mock_task_service.go -package=mocks

// TaskServiceInterface define el contrato para el servicio de tareas
type TaskServiceInterface interface {
	// CreateTask crea una nueva tarea
	CreateTask(ctx context.Context, title, description string) (*domain.Task, error)
	
	// GetTaskByID obtiene una tarea por su ID
	GetTaskByID(ctx context.Context, id int) (*domain.Task, error)
	
	// GetAllTasks obtiene todas las tareas
	GetAllTasks(ctx context.Context) ([]*domain.Task, error)
	
	// UpdateTask actualiza una tarea existente
	UpdateTask(ctx context.Context, id int, title, description string, completed *bool) (*domain.Task, error)
	
	// DeleteTask elimina una tarea por su ID
	DeleteTask(ctx context.Context, id int) error
	
	// GetTasksByStatus obtiene tareas filtradas por estado de completado
	GetTasksByStatus(ctx context.Context, completed bool) ([]*domain.Task, error)
	
	// MarkTaskAsCompleted marca una tarea como completada
	MarkTaskAsCompleted(ctx context.Context, id int) (*domain.Task, error)
	
	// MarkTaskAsUncompleted marca una tarea como no completada
	MarkTaskAsUncompleted(ctx context.Context, id int) (*domain.Task, error)
}