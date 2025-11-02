package application

import (
	"context"
	"fmt"

	"github.com/YerkoTenorio/api-go-hexagonal/modules/task/domain"
)

// TaskService maneja los casos de uso relacionados con tareas
type TaskService struct {
	taskRepo domain.TaskRepository
}

// NewTaskService crea una nueva instancia de TaskService
func NewTaskService(taskRepo domain.TaskRepository) *TaskService {
	return &TaskService{
		taskRepo: taskRepo,
	}
}

// CreateTask crea una nueva tarea
func (s *TaskService) CreateTask(ctx context.Context, title, description string) (*domain.Task, error) {
	// Validación
	if title == "" {
		return nil, fmt.Errorf("el titulo es requerido")
	}
	if description == "" {
		return nil, fmt.Errorf("la descripcion es requerida")
	}

	// Crear nueva tarea
	task := domain.NewTask(title, description)

	if !task.IsValid() {
		return nil, fmt.Errorf("la tarea no es válida")
	}

	// Persistir usando el repositorio
	return s.taskRepo.Create(ctx, task)
}

// GetTaskByID obtiene una tarea por su ID
func (s *TaskService) GetTaskByID(ctx context.Context, id int) (*domain.Task, error) {
	if id == 0 {
		return nil, fmt.Errorf("el ID de la tarea no puede ser cero")
	}

	task, err := s.taskRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("no se pudo obtener la tarea con ID %d: %w", id, err)
	}

	return task, nil
}

// GetAllTasks obtiene todas las tareas
func (s *TaskService) GetAllTasks(ctx context.Context) ([]*domain.Task, error) {
	tasks, err := s.taskRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("no se pudieron obtener las tareas: %w", err)
	}

	return tasks, nil
}

// UpdateTask actualiza una tarea existente
func (s *TaskService) UpdateTask(ctx context.Context, id int, title, description string, completed *bool) (*domain.Task, error) {
	if id == 0 {
		return nil, fmt.Errorf("El ID de la tarea es requerido")
	}

	// Obtener la tarea existente
	task, err := s.taskRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("no se pudo encontrar la tarea con ID %d: %w", id, err)
	}

	// Preparar los valores finales para la actualización
	finalTitle := task.Title             // Valor actual por defecto
	finalDescription := task.Description // Valor actual por defecto

	// Actualizar solo los campos proporcionados (no vacíos)
	if title != "" {
		finalTitle = title
	}
	if description != "" {
		finalDescription = description
	}

	// Aplicar la actualización UNA SOLA VEZ con los valores finales
	task.Update(finalTitle, finalDescription)

	// Actualizar estado de completado si se proporciona
	if completed != nil {
		if *completed {
			task.MarkAsCompleted()
		} else {
			task.MarkAsUncompleted()
		}
	}

	// Validar la tarea actualizada
	if !task.IsValid() {
		return nil, fmt.Errorf("la tarea actualizada no es valida")
	}

	// AQUÍ ES DONDE SE GUARDAN LOS CAMBIOS EN LA BASE DE DATOS
	updatedTask, err := s.taskRepo.Update(ctx, task)
	if err != nil {
		return nil, fmt.Errorf("no se pudo actualizar la tarea: %w", err)
	}

	return updatedTask, nil
}

// DeleteTask elimina una tarea por su ID
func (s *TaskService) DeleteTask(ctx context.Context, id int) error {
	if id == 0 {
		return fmt.Errorf("el ID de la tarea es requerido")
	}

	// Verificar que la tarea existe antes de eliminarla
	_, err := s.taskRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("no se pudo encontrar la tarea con ID %d: %w", id, err)
	}

	// Eliminar la tarea
	err = s.taskRepo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("no se pudo eliminar la tarea con ID %d: %w", id, err)
	}

	return nil
}

// GetTasksByStatus obtiene tareas filtradas por estado de completado
func (s *TaskService) GetTasksByStatus(ctx context.Context, completed bool) ([]*domain.Task, error) {
	tasks, err := s.taskRepo.GetByStatus(ctx, completed)
	if err != nil {
		return nil, fmt.Errorf("no se pudieron obtener las tareas con estado completado=%t: %w", completed, err)
	}

	return tasks, nil
}

// MarkTaskAsCompleted marca una tarea como completada
func (s *TaskService) MarkTaskAsCompleted(ctx context.Context, id int) (*domain.Task, error) {
	if id == 0 {
		return nil, fmt.Errorf("el ID de la tarea es requerido")
	}

	// Obtener la tarea existente
	task, err := s.taskRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("no se pudo encontrar la tarea con ID %d: %w", id, err)
	}

	// Marcar como completada
	task.MarkAsCompleted()

	// Persistir los cambios
	updatedTask, err := s.taskRepo.Update(ctx, task)
	if err != nil {
		return nil, fmt.Errorf("no se pudo marcar la tarea como completada: %w", err)
	}

	return updatedTask, nil
}

// MarkTaskAsUncompleted marca una tarea como no completada
func (s *TaskService) MarkTaskAsUncompleted(ctx context.Context, id int) (*domain.Task, error) {
	if id == 0 {
		return nil, fmt.Errorf("el ID de la tarea es requerido")
	}

	// Obtener la tarea existente
	task, err := s.taskRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("no se pudo encontrar la tarea con ID %d: %w", id, err)
	}

	// Marcar como no completada
	task.MarkAsUncompleted()

	// Persistir los cambios
	updatedTask, err := s.taskRepo.Update(ctx, task)
	if err != nil {
		return nil, fmt.Errorf("no se pudo marcar la tarea como no completada: %w", err)
	}

	return updatedTask, nil
}
