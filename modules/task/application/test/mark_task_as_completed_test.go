package application_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/YerkoTenorio/api-go-hexagonal/modules/task/application"
	"github.com/YerkoTenorio/api-go-hexagonal/modules/task/application/mocks"
	"github.com/YerkoTenorio/api-go-hexagonal/modules/task/domain"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

// TestTaskService_MarkTaskAsCompleted_Success verifica que se puede marcar una tarea como completada
func TestTaskService_MarkTaskAsCompleted_Success(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTaskRepository(ctrl)
	service := application.NewTaskService(mockRepo)

	// Tarea existente (incompleta)
	existingTask := &domain.Task{
		ID:          1,
		Title:       "Tarea a completar",
		Description: "Descripción de la tarea",
		Completed:   false,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	// Tarea después de marcarla como completada
	completedTask := &domain.Task{
		ID:          1,
		Title:       "Tarea a completar",
		Description: "Descripción de la tarea",
		Completed:   true,
		CreatedAt:   existingTask.CreatedAt,
		UpdatedAt:   time.Now().UTC(),
	}

	// Expectativas del mock
	mockRepo.EXPECT().
		GetByID(gomock.Any(), 1).
		Return(existingTask, nil).
		Times(1)

	mockRepo.EXPECT().
		Update(gomock.Any(), gomock.Any()).
		Return(completedTask, nil).
		Times(1)

	// Act
	result, err := service.MarkTaskAsCompleted(context.Background(), 1)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, result.ID)
	assert.True(t, result.Completed)
	assert.Equal(t, "Tarea a completar", result.Title)
}

// TestTaskService_MarkTaskAsCompleted_InvalidID_ShouldReturnError verifica el manejo de ID inválido
func TestTaskService_MarkTaskAsCompleted_InvalidID_ShouldReturnError(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTaskRepository(ctrl)
	service := application.NewTaskService(mockRepo)

	// Act
	result, err := service.MarkTaskAsCompleted(context.Background(), 0)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "el ID de la tarea es requerido")
}

// TestTaskService_MarkTaskAsCompleted_TaskNotFound_ShouldReturnError verifica el manejo cuando la tarea no existe
func TestTaskService_MarkTaskAsCompleted_TaskNotFound_ShouldReturnError(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTaskRepository(ctrl)
	service := application.NewTaskService(mockRepo)

	notFoundError := errors.New("task not found")

	// Expectativa del mock
	mockRepo.EXPECT().
		GetByID(gomock.Any(), 999).
		Return(nil, notFoundError).
		Times(1)

	// Act
	result, err := service.MarkTaskAsCompleted(context.Background(), 999)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "no se pudo encontrar la tarea con ID 999")
	assert.Contains(t, err.Error(), "task not found")
}

// TestTaskService_MarkTaskAsCompleted_UpdateError_ShouldReturnError verifica el manejo de errores en la actualización
func TestTaskService_MarkTaskAsCompleted_UpdateError_ShouldReturnError(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTaskRepository(ctrl)
	service := application.NewTaskService(mockRepo)

	existingTask := &domain.Task{
		ID:          1,
		Title:       "Tarea a completar",
		Description: "Descripción de la tarea",
		Completed:   false,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	updateError := errors.New("database update failed")

	// Expectativas del mock
	mockRepo.EXPECT().
		GetByID(gomock.Any(), 1).
		Return(existingTask, nil).
		Times(1)

	mockRepo.EXPECT().
		Update(gomock.Any(), gomock.Any()).
		Return(nil, updateError).
		Times(1)

	// Act
	result, err := service.MarkTaskAsCompleted(context.Background(), 1)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "no se pudo marcar la tarea como completada")
	assert.Contains(t, err.Error(), "database update failed")
}

// TestTaskService_MarkTaskAsCompleted_AlreadyCompleted_ShouldStillWork verifica que funciona aunque ya esté completada
func TestTaskService_MarkTaskAsCompleted_AlreadyCompleted_ShouldStillWork(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTaskRepository(ctrl)
	service := application.NewTaskService(mockRepo)

	// Tarea ya completada
	alreadyCompletedTask := &domain.Task{
		ID:          1,
		Title:       "Tarea ya completada",
		Description: "Descripción de la tarea",
		Completed:   true,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	// Expectativas del mock
	mockRepo.EXPECT().
		GetByID(gomock.Any(), 1).
		Return(alreadyCompletedTask, nil).
		Times(1)

	mockRepo.EXPECT().
		Update(gomock.Any(), gomock.Any()).
		Return(alreadyCompletedTask, nil).
		Times(1)

	// Act
	result, err := service.MarkTaskAsCompleted(context.Background(), 1)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.Completed)
}
