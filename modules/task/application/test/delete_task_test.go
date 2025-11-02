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

// TestTaskService_DeleteTask_Success verifica que se puede eliminar una tarea existente
func TestTaskService_DeleteTask_Success(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish() // Limpiar el controlador al finalizar el test

	mockRepo := mocks.NewMockTaskRepository(ctrl)
	service := application.NewTaskService(mockRepo)

	existingTask := &domain.Task{
		ID:          1,
		Title:       "Tarea a eliminar",
		Description: "Descripcion de la tarea",
		Completed:   false,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	// Expectativas del mock

	mockRepo.EXPECT().
		GetByID(gomock.Any(), 1).
		Return(existingTask, nil).
		Times(1)

	mockRepo.EXPECT().
		Delete(gomock.Any(), 1).
		Return(nil).
		Times(1)

	// Act
	err := service.DeleteTask(context.Background(), 1)
	// Assert
	assert.NoError(t, err) // Verificar que no se retorna un error

}

// TestTaskService_DeleteTask_ZeroID_ShouldReturnError verifica que se retorna un error cuando se intenta eliminar una tarea con ID cero
func TestTaskService_DeleteTask_ZeroID_ShouldReturnError(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTaskRepository(ctrl)
	service := application.NewTaskService(mockRepo)

	// No esperamos llamadas al repositorio porque la validación falla antes

	// Act
	err := service.DeleteTask(context.Background(), 0)

	// Assert
	assert.Error(t, err)                                           // Verificar que se retorna un error
	assert.Equal(t, "el ID de la tarea es requerido", err.Error()) // Verificar que el error es el esperado
}

func TestTaskService_DeleteTask_RepositoryDeleteError_ShouldPropagateError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTaskRepository(ctrl)
	service := application.NewTaskService(mockRepo)

	existingTask := &domain.Task{
		ID:          1,
		Title:       "Tarea existente",
		Description: "Descripcion",
		Completed:   false,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	// Expectativas del mock
	mockRepo.EXPECT().
		GetByID(gomock.Any(), 1).
		Return(existingTask, nil).
		Times(1)

	mockRepo.EXPECT().
		Delete(gomock.Any(), 1).
		Return(errors.New("database error")).
		Times(1)

	// Act
	err := service.DeleteTask(context.Background(), 1)

	// Assert
	assert.Error(t, err)                                                   // Verificar que se retorna un error
	assert.Contains(t, err.Error(), "no se pudo eliminar la tarea con ID") // Verificar que el error contiene el mensaje esperado
	assert.Contains(t, err.Error(), "database error")                      // Verificar que el error contiene el mensaje esperado

}
