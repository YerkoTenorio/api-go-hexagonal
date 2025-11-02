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

// TestTaskService_UpdateTask_Success_AllFields verifica que se puede actualizar una tarea con todos los campos
func TestTaskService_UpdateTask_Success_AllFields(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTaskRepository(ctrl)
	service := application.NewTaskService(mockRepo)

	// datos de prueba
	existingTask := &domain.Task{
		ID:          1,
		Title:       "Titulo valido",
		Description: "Descripcion valida",
		Completed:   false,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	updatedTask := &domain.Task{
		ID:          1,
		Title:       "Nuevo titulo",
		Description: "Nueva descripcion",
		Completed:   true,
		CreatedAt:   existingTask.CreatedAt,
		UpdatedAt:   time.Now().UTC(),
	}

	completed := true

	// Expectativas del mock
	mockRepo.EXPECT().
		GetByID(gomock.Any(), 1).
		Return(existingTask, nil).
		Times(1)

	mockRepo.EXPECT().
		Update(gomock.Any(), gomock.Any()).
		Return(updatedTask, nil).
		Times(1)

	// Act
	result, err := service.UpdateTask(context.Background(), 1, "New Title", "New Description", &completed)

	// Assert
	assert.NoError(t, err)                                   // Verificar que no se retorna un error
	assert.NotNil(t, result)                                 // Verificar que se retorna una tarea no nula
	assert.Equal(t, "Nuevo titulo", result.Title)            // Verificar que el titulo de la tarea es el esperado
	assert.Equal(t, "Nueva descripcion", result.Description) // Verificar que la descripcion de la tarea es la esperada
	assert.True(t, result.Completed)                         // Verificar que la tarea esta completada
}

// TestTaskService_UpdateTask_Success_OnlyTitle verifica que se puede actualizar una tarea con solo el titulo
func TestTaskService_UpdateTask_Success_OnlyTitle(t *testing.T) {
	//Arrange

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTaskRepository(ctrl)
	service := application.NewTaskService(mockRepo)

	existingTask := &domain.Task{
		ID:          1,
		Title:       "Old Title",
		Description: "Keep Description",
		Completed:   false,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	updatedTask := &domain.Task{
		ID:          1,
		Title:       "New Title",
		Description: "Keep Description",
		Completed:   false,
		CreatedAt:   existingTask.CreatedAt,
		UpdatedAt:   time.Now().UTC(),
	}

	mockRepo.EXPECT().
		GetByID(gomock.Any(), 1).
		Return(existingTask, nil).
		Times(1)

	mockRepo.EXPECT().
		Update(gomock.Any(), gomock.Any()).
		Return(updatedTask, nil).
		Times(1)

	// Act
	result, err := service.UpdateTask(context.Background(), 1, "New Title", "", nil)

	// Assert
	assert.NoError(t, err)                                  // Verificar que no se retorna un error
	assert.NotNil(t, result)                                // Verificar que se retorna una tarea no nula
	assert.Equal(t, "New Title", result.Title)              // Verificar que el titulo de la tarea es el esperado
	assert.Equal(t, "Keep Description", result.Description) // Verificar que la descripcion de la tarea es la esperada
	assert.False(t, result.Completed)                       // Verificar que la tarea no esta completada
}

// TestTaskService_UpdateTask_ZeroID_ShouldReturnError verifica que se retorna un error cuando se actualiza una tarea con ID 0

func TestTaskService_UpdateTask_ZeroID_ShouldReturnError(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTaskRepository(ctrl)
	service := application.NewTaskService(mockRepo)

	// No esperamos llamadas al repositorio

	// Act

	result, err := service.UpdateTask(context.Background(), 0, "Title", "Description", nil)

	// Assert
	assert.Error(t, err)                                           // Verificar que se retorna un error
	assert.Nil(t, result)                                          // Verificar que no se retorna una tarea
	assert.Equal(t, "El ID de la tarea es requerido", err.Error()) // Verificar que el error es el esperado
}

func TestTaskService_UpdateTask_RepositoryUpdateError_ShouldPropagateError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTaskRepository(ctrl)
	service := application.NewTaskService(mockRepo)

	existingTask := &domain.Task{
		ID:          1,
		Title:       "Title",
		Description: "Description",
		Completed:   false,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	updateError := errors.New("database update failed")

	mockRepo.EXPECT().
		GetByID(gomock.Any(), 1).
		Return(existingTask, nil).
		Times(1)

	mockRepo.EXPECT().
		Update(gomock.Any(), gomock.Any()).
		Return(nil, updateError).
		Times(1)

	// Act
	result, err := service.UpdateTask(context.Background(), 1, "New Title", "New Description", nil)

	// Assert
	assert.Error(t, err)  // Verificar que se retorna un error
	assert.Nil(t, result) // Verificar que no se retorna una tarea
	assert.Contains(t, err.Error(), "no se pudo actualizar la tarea")
	assert.Contains(t, err.Error(), "database update failed")
}
