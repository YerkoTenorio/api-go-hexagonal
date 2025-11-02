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

// TestTaskService_GetTasksByStatus_Success_CompletedTasks verifica que se pueden obtener tareas completadas
func TestTaskService_GetTasksByStatus_Success_CompletedTasks(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTaskRepository(ctrl)
	service := application.NewTaskService(mockRepo)

	// Tareas completadas de prueba
	completedTasks := []*domain.Task{
		{
			ID:          1,
			Title:       "Tarea completada 1",
			Description: "Descripción 1",
			Completed:   true,
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
		},
		{
			ID:          2,
			Title:       "Tarea completada 2",
			Description: "Descripción 2",
			Completed:   true,
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
		},
	}

	// Expectativa del mock
	mockRepo.EXPECT().
		GetByStatus(gomock.Any(), true).
		Return(completedTasks, nil).
		Times(1)

	// Act
	result, err := service.GetTasksByStatus(context.Background(), true)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	assert.Equal(t, completedTasks[0].ID, result[0].ID)
	assert.Equal(t, completedTasks[1].ID, result[1].ID)
	assert.True(t, result[0].Completed)
	assert.True(t, result[1].Completed)
}

// TestTaskService_GetTasksByStatus_Success_IncompleteTasks verifica que se pueden obtener tareas incompletas
func TestTaskService_GetTasksByStatus_Success_IncompleteTasks(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTaskRepository(ctrl)
	service := application.NewTaskService(mockRepo)

	// Tareas incompletas de prueba
	incompleteTasks := []*domain.Task{
		{
			ID:          3,
			Title:       "Tarea pendiente 1",
			Description: "Descripción 3",
			Completed:   false,
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
		},
		{
			ID:          4,
			Title:       "Tarea pendiente 2",
			Description: "Descripción 4",
			Completed:   false,
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
		},
	}

	// Expectativa del mock
	mockRepo.EXPECT().
		GetByStatus(gomock.Any(), false).
		Return(incompleteTasks, nil).
		Times(1)

	// Act
	result, err := service.GetTasksByStatus(context.Background(), false)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	assert.Equal(t, incompleteTasks[0].ID, result[0].ID)
	assert.Equal(t, incompleteTasks[1].ID, result[1].ID)
	assert.False(t, result[0].Completed)
	assert.False(t, result[1].Completed)
}

// TestTaskService_GetTasksByStatus_Success_EmptyResult verifica el comportamiento con resultado vacío
func TestTaskService_GetTasksByStatus_Success_EmptyResult(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTaskRepository(ctrl)
	service := application.NewTaskService(mockRepo)

	// Lista vacía
	emptyTasks := []*domain.Task{}

	// Expectativa del mock
	mockRepo.EXPECT().
		GetByStatus(gomock.Any(), true).
		Return(emptyTasks, nil).
		Times(1)

	// Act
	result, err := service.GetTasksByStatus(context.Background(), true)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 0)
}

// TestTaskService_GetTasksByStatus_RepositoryError_ShouldPropagateError verifica el manejo de errores del repositorio
func TestTaskService_GetTasksByStatus_RepositoryError_ShouldPropagateError(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTaskRepository(ctrl)
	service := application.NewTaskService(mockRepo)

	repositoryError := errors.New("database connection failed")

	// Expectativa del mock
	mockRepo.EXPECT().
		GetByStatus(gomock.Any(), true).
		Return(nil, repositoryError).
		Times(1)

	// Act
	result, err := service.GetTasksByStatus(context.Background(), true)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "no se pudieron obtener las tareas con estado completado=true")
	assert.Contains(t, err.Error(), "database connection failed")
}
