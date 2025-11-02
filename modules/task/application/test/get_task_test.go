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

// TestTaskService_GetTaskByID_Success verifica que se puede obtener una tarea por ID exitosamente
func TestTaskService_GetTaskByID_Success(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTaskRepository(ctrl)
	service := application.NewTaskService(mockRepo)

	expectedTask := &domain.Task{
		ID:          1,
		Title:       "Titulo valido",
		Description: "Descripcion valida",
		Completed:   false,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	mockRepo.EXPECT().
		GetByID(gomock.Any(), 1).
		Return(expectedTask, nil).
		Times(1)

	// Act
	result, err := service.GetTaskByID(context.Background(), 1)

	// asset

	assert.NoError(t, err)                                        // Verificar que no se retorna un error
	assert.NotNil(t, result)                                      // Verificar que se retorna una tarea
	assert.Equal(t, expectedTask.ID, result.ID)                   // Verificar que el ID de la tarea es el esperado
	assert.Equal(t, expectedTask.Title, result.Title)             // Verificar que el titulo de la tarea es el esperado
	assert.Equal(t, expectedTask.Description, result.Description) // Verificar que la descripcion de la tarea es la esperada
	assert.Equal(t, expectedTask.Completed, result.Completed)     // Verificar que el estado de completado de la tarea es el esperado

}

// TestTaskService_GetTaskByID_ZeroID_ShouldReturnError verifica que se retorna un error cuando se solicita una tarea con ID cero
func TestTaskService_GetTaskByID_ZeroID_ShouldReturnError(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTaskRepository(ctrl)
	service := application.NewTaskService(mockRepo)

	// No esperamos llamados al reposotprop poqie falla antes de llegar ahi

	// Act
	result, err := service.GetTaskByID(context.Background(), 0)

	// Assert
	assert.Error(t, err)                                                // Verificar que se retorna un error
	assert.Nil(t, result)                                               // Verificar que no se retorna una tarea
	assert.Equal(t, "el ID de la tarea no puede ser cero", err.Error()) // Verificar que el error contiene el mensaje esperado
}

// TestTaskService_GetTaskByID_RepositoryError_ShouldPropagateError verifica que se propaga un error del repositorio cuando se solicita una tarea por ID
func TestTaskService_GetTaskByID_RepositoryError_ShouldPropagateError(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTaskRepository(ctrl)
	service := application.NewTaskService(mockRepo)

	repositoryError := errors.New("database connection error")
	taskID := 1

	mockRepo.EXPECT().
		GetByID(gomock.Any(), taskID).
		Return(nil, repositoryError).
		Times(1)

	// Act
	result, err := service.GetTaskByID(context.Background(), taskID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "no se pudo obtener la tarea con ID 1")
	assert.Contains(t, err.Error(), "database connection error")
}

func TestTaskService_GetTaskByID_TaskNotFound_ShouldReturnError(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTaskRepository(ctrl)
	service := application.NewTaskService(mockRepo)

	notFoundError := errors.New("task not found")
	taskID := 999

	mockRepo.EXPECT().
		GetByID(gomock.Any(), taskID).
		Return(nil, notFoundError).
		Times(1)

	// Act
	result, err := service.GetTaskByID(context.Background(), taskID)

	// Assert
	assert.Error(t, err)  // Verificar que se retorna un error
	assert.Nil(t, result) // Verificar que no se retorna una tarea
	assert.Contains(t, err.Error(), "no se pudo obtener la tarea con ID 999")
	assert.Contains(t, err.Error(), "task not found")
}

// tests para GetAllTasks
// TestTaskService_GetAllTasks_Success verifica que se puede obtener todas las tareas exitosamente
func TestTaskService_GetAllTasks_Success(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTaskRepository(ctrl)
	service := application.NewTaskService(mockRepo)

	expectedTasks := []*domain.Task{
		{
			ID:          1,
			Title:       "Titulo valido",
			Description: "Descripcion valida",
			Completed:   false,
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
		},
		{
			ID:          2,
			Title:       "Titulo valido 2",
			Description: "Descripcion valida 2",
			Completed:   true,
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
		},
		{
			ID:          3,
			Title:       "Titulo valido 3",
			Description: "Descripcion valida 3",
			Completed:   false,
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
		},
	}

	mockRepo.EXPECT().
		GetAll(gomock.Any()).
		Return(expectedTasks, nil).
		Times(1)

	// Act
	result, err := service.GetAllTasks(context.Background())

	// Assert
	assert.NoError(t, err)                             //Verificar que no hay error
	assert.NotNil(t, result)                           //Verificar que se retorna una lista de tareas
	assert.Len(t, result, 3)                           //Verificar que se retorna 3 tareas
	assert.Equal(t, expectedTasks[0].ID, result[0].ID) //Verificar que el ID de la tarea es el esperado
	assert.Equal(t, expectedTasks[1].ID, result[1].ID) // Verificar que el ID de la tarea es el esperado
	assert.Equal(t, expectedTasks[2].ID, result[2].ID) // Verificar que el ID de la tarea es el esperado

}

// TestTaskService_GetAllTasks_EmptyResult_ShouldReturnEmptySlice verifica que se retorna un slice vacio cuando no hay tareas disponibles
func TestTaskService_GetAllTasks_EmptyResult_ShouldReturnEmptySlice(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTaskRepository(ctrl)
	service := application.NewTaskService(mockRepo)

	emptyTasks := []*domain.Task{}

	mockRepo.EXPECT().
		GetAll(gomock.Any()).
		Return(emptyTasks, nil).
		Times(1)

	// Act
	result, err := service.GetAllTasks(context.Background())

	// Assert
	assert.NoError(t, err)   //Verificar que no hay error
	assert.NotNil(t, result) //Verificar que se retorna una lista de tareas
	assert.Len(t, result, 0) //Verificar que se retorna 0 tareas
}

func TestTaskService_GetAllTasks_RepositoryError_ShouldPropagateError(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTaskRepository(ctrl)
	service := application.NewTaskService(mockRepo)

	repositoryError := errors.New("database error")

	mockRepo.EXPECT().
		GetAll(gomock.Any()).
		Return(nil, repositoryError).
		Times(1)

	// Act
	result, err := service.GetAllTasks(context.Background())

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "database error")
}
