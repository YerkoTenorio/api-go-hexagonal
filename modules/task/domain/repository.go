package domain

import (
	"context"
)

//go:generate mockgen -source=repository.go -destination=../application/mocks/mock_task_repository.go -package=mocks

// TaskRepository define el contrato para el repositorio de tareas

type TaskRepository interface {
	// Create guarda una nueva tarea en el repositorio
	Create(ctx context.Context, task *Task) (*Task, error)
	// GetByID obtiene una tarea por si ID
	GetByID(ctx context.Context, id int) (*Task, error)
	// GetAll obtiene todas las tareas
	GetAll(ctx context.Context) ([]*Task, error)
	// Update actualiza una tarea por su id
	Update(ctx context.Context, task *Task) (*Task, error)
	// Delete elimina una tarea por su id
	Delete(ctx context.Context, id int) error
	// GetByStatus obtiene tareas por su estado
	GetByStatus(ctx context.Context, completed bool) ([]*Task, error)
}
