package domain

import "time"

// Task representa una tarea en el sistema
type Task struct {
	ID          int       `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	Completed   bool      `json:"completed" db:"completed"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// NewTask crea una nueva instancia de Task
func NewTask(title, description string) *Task {
	// agregar el tiempo en UTC
	now := time.Now().UTC()
	return &Task{
		Title:       title,
		Description: description,
		Completed:   false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// MarkAsCompleted marca la tarea como completada
func (t *Task) MarkAsCompleted() {
	t.Completed = true
	t.UpdatedAt = time.Now().UTC()
}

// MarkAsUncompleted marca la tarea como incompleta
func (t *Task) MarkAsUncompleted() {
	t.Completed = false
	t.UpdatedAt = time.Now().UTC()
}

// Update actualiza los campos de la tarea
func (t *Task) Update(title, description string) {
	if title != "" {
		t.Title = title
	}
	if description != "" {
		t.Description = description
	}
	t.UpdatedAt = time.Now().UTC()
}

// IsValid valida que la tarea tenga los campos requeridos
func (t *Task) IsValid() bool {
	return t.Title != "" && t.Description != ""
}
