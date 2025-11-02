package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestNewTask verifica la creación correcta de una nueva tarea
func TestNewTask(t *testing.T) {
	// Arrange
	title := "Test Task"
	description := "Test Description"
	beforeCreation := time.Now().UTC()

	// Act
	task := NewTask(title, description)

	// Assert
	assert.Equal(t, title, task.Title)
	assert.Equal(t, description, task.Description)
	assert.False(t, task.Completed)
	assert.Equal(t, 0, task.ID) // ID debe ser 0 por defecto

	// Verificar que los timestamps están cerca del momento de creación
	assert.True(t, task.CreatedAt.After(beforeCreation) || task.CreatedAt.Equal(beforeCreation))
	assert.True(t, task.UpdatedAt.After(beforeCreation) || task.UpdatedAt.Equal(beforeCreation))
	assert.Equal(t, task.CreatedAt, task.UpdatedAt) // Deben ser iguales al crear
}

// TestNewTask_EmptyFields verifica la creación con campos vacíos
func TestNewTask_EmptyFields(t *testing.T) {
	// Act
	task := NewTask("", "")

	// Assert
	assert.Equal(t, "", task.Title)
	assert.Equal(t, "", task.Description)
	assert.False(t, task.Completed)
	assert.NotZero(t, task.CreatedAt)
	assert.NotZero(t, task.UpdatedAt)
}

// TestTask_MarkAsCompleted verifica marcar tarea como completada
func TestTask_MarkAsCompleted(t *testing.T) {
	// Arrange
	task := NewTask("Test", "Description")
	originalUpdatedAt := task.UpdatedAt

	// Esperar un poco para asegurar diferencia en timestamp
	time.Sleep(1 * time.Millisecond)

	// Act
	task.MarkAsCompleted()

	// Assert
	assert.True(t, task.Completed)
	assert.True(t, task.UpdatedAt.After(originalUpdatedAt))
}

// TestTask_MarkAsUncompleted verifica marcar tarea como incompleta
func TestTask_MarkAsUncompleted(t *testing.T) {
	// Arrange
	task := NewTask("Test", "Description")
	task.MarkAsCompleted() // Primero la marcamos como completada
	originalUpdatedAt := task.UpdatedAt

	// Esperar un poco para asegurar diferencia en timestamp
	time.Sleep(1 * time.Millisecond)

	// Act
	task.MarkAsUncompleted()

	// Assert
	assert.False(t, task.Completed)
	assert.True(t, task.UpdatedAt.After(originalUpdatedAt))
}

// TestTask_Update_BothFields verifica actualización de ambos campos
func TestTask_Update_BothFields(t *testing.T) {
	// Arrange
	task := NewTask("Original Title", "Original Description")
	originalUpdatedAt := task.UpdatedAt
	newTitle := "Updated Title"
	newDescription := "Updated Description"

	// Esperar un poco para asegurar diferencia en timestamp
	time.Sleep(1 * time.Millisecond)

	// Act
	task.Update(newTitle, newDescription)

	// Assert
	assert.Equal(t, newTitle, task.Title)
	assert.Equal(t, newDescription, task.Description)
	assert.True(t, task.UpdatedAt.After(originalUpdatedAt))
}

// TestTask_Update_TitleOnly verifica actualización solo del título
func TestTask_Update_TitleOnly(t *testing.T) {
	// Arrange
	task := NewTask("Original Title", "Original Description")
	originalDescription := task.Description
	originalUpdatedAt := task.UpdatedAt
	newTitle := "Updated Title"

	// Esperar un poco para asegurar diferencia en timestamp
	time.Sleep(1 * time.Millisecond)

	// Act
	task.Update(newTitle, "")

	// Assert
	assert.Equal(t, newTitle, task.Title)
	assert.Equal(t, originalDescription, task.Description) // No debe cambiar
	assert.True(t, task.UpdatedAt.After(originalUpdatedAt))
}

// TestTask_Update_DescriptionOnly verifica actualización solo de la descripción
func TestTask_Update_DescriptionOnly(t *testing.T) {
	// Arrange
	task := NewTask("Original Title", "Original Description")
	originalTitle := task.Title
	originalUpdatedAt := task.UpdatedAt
	newDescription := "Updated Description"

	// Esperar un poco para asegurar diferencia en timestamp
	time.Sleep(1 * time.Millisecond)

	// Act
	task.Update("", newDescription)

	// Assert
	assert.Equal(t, originalTitle, task.Title) // No debe cambiar
	assert.Equal(t, newDescription, task.Description)
	assert.True(t, task.UpdatedAt.After(originalUpdatedAt))
}

// TestTask_Update_EmptyFields verifica que campos vacíos no actualizan
func TestTask_Update_EmptyFields(t *testing.T) {
	// Arrange
	task := NewTask("Original Title", "Original Description")
	originalTitle := task.Title
	originalDescription := task.Description
	originalUpdatedAt := task.UpdatedAt

	// Esperar un poco para asegurar diferencia en timestamp
	time.Sleep(1 * time.Millisecond)

	// Act
	task.Update("", "")

	// Assert
	assert.Equal(t, originalTitle, task.Title)              // No debe cambiar
	assert.Equal(t, originalDescription, task.Description)  // No debe cambiar
	assert.True(t, task.UpdatedAt.After(originalUpdatedAt)) // UpdatedAt sí se actualiza
}

// TestTask_IsValid_ValidTask verifica validación de tarea válida
func TestTask_IsValid_ValidTask(t *testing.T) {
	// Arrange
	task := NewTask("Valid Title", "Valid Description")

	// Act
	isValid := task.IsValid()

	// Assert
	assert.True(t, isValid)
}

// TestTask_IsValid_EmptyTitle verifica validación con título vacío
func TestTask_IsValid_EmptyTitle(t *testing.T) {
	// Arrange
	task := NewTask("", "Valid Description")

	// Act
	isValid := task.IsValid()

	// Assert
	assert.False(t, isValid)
}

// TestTask_IsValid_EmptyDescription verifica validación con descripción vacía
func TestTask_IsValid_EmptyDescription(t *testing.T) {
	// Arrange
	task := NewTask("Valid Title", "")

	// Act
	isValid := task.IsValid()

	// Assert
	assert.False(t, isValid)
}

// TestTask_IsValid_BothEmpty verifica validación con ambos campos vacíos
func TestTask_IsValid_BothEmpty(t *testing.T) {
	// Arrange
	task := NewTask("", "")

	// Act
	isValid := task.IsValid()

	// Assert
	assert.False(t, isValid)
}

// TestTask_IsValid_WhitespaceFields verifica validación con espacios en blanco
func TestTask_IsValid_WhitespaceFields(t *testing.T) {
	// Arrange
	task := NewTask("   ", "   ")

	// Act
	isValid := task.IsValid()

	// Assert
	// Nota: La implementación actual considera espacios como válidos
	// Si quisieras cambiar esto, tendrías que usar strings.TrimSpace()
	assert.True(t, isValid)
}
