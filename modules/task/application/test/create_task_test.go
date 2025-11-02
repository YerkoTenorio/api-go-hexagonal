package application_test

import (
	"context"
	"testing"
	"time"

	"github.com/YerkoTenorio/api-go-hexagonal/modules/task/application"
	"github.com/YerkoTenorio/api-go-hexagonal/modules/task/application/mocks"
	"github.com/YerkoTenorio/api-go-hexagonal/modules/task/domain"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

// TestTaskService_CreateTask_Success verifica que se puede crear una tarea exitosamente
func TestTaskService_CreateTask_Success(t *testing.T) {
	// Arrange (Preparar)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Crear un mock del repositorio
	mockRepo := mocks.NewMockTaskRepository(ctrl)
	service := application.NewTaskService(mockRepo)
	// Preparar el contexto
	ctx := context.Background()
	title := "Tarea de prueba"
	description := "Descripcion de prueba"

	// Crear la tarea esperada que devolveria el repositorio
	expectedTask := &domain.Task{
		ID:          1,
		Title:       title,
		Description: description,
		Completed:   false,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}
	// Configurar expectativa del mock
	mockRepo.EXPECT().Create(ctx, gomock.Any()).Return(expectedTask, nil).Times(1)
	// Act (Actuar)
	result, err := service.CreateTask(ctx, title, description)

	// Assert (Verificar)

	assert.NoError(t, err)                           // Verificar que no hay error
	assert.NotNil(t, result)                         // Verificar que se retorna una tarea no nula
	assert.Equal(t, expectedTask.ID, result.ID)      // Verificar que el ID de la tarea es el esperado
	assert.Equal(t, title, result.Title)             // Verificar que el titulo de la tarea es el esperado
	assert.Equal(t, description, result.Description) // Verificar que la descripcion de la tarea es la esperada
	assert.False(t, result.Completed)                // Verificar que la tarea no esta completada
	assert.NotZero(t, result.CreatedAt)              // Verificar que la fecha de creacion no es cero
	assert.NotZero(t, result.UpdatedAt)              // Verificar que la fecha de actualizacion no es cero
}

// TestTaskService_CreateTask_EmptyTitle_ShouldReturnError verifica que no se puede crear una tarea con un titulo vacio
func TestTaskService_CreateTask_EmptyTitle_ShouldReturnError(t *testing.T) {
	// Arrange (Preparar)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTaskRepository(ctrl)
	service := application.NewTaskService(mockRepo)

	ctx := context.Background()
	title := ""
	description := "Descripcion valida"

	// no configuramos expectativas en el mock porque la validacion debe fallar antes de llamar al repositorio

	result, err := service.CreateTask(ctx, title, description)

	//Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "el titulo es requerido")

}

// TestTaskService_CreateTask_EmptyDescription_ShouldReturnError verifica que no se puede crear una tarea con una descripcion vacia
func TestTaskService_CreateTask_EmptyDescription_ShouldReturnError(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTaskRepository(ctrl)
	service := application.NewTaskService(mockRepo)

	ctx := context.Background()
	title := "Titulo valido"
	description := ""

	// Act
	result, err := service.CreateTask(ctx, title, description)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "la descripcion es requerida") // Verificar que el error contiene el mensaje esperado

}

// TestTaskService_CreateTask_RepositoryError_ShouldPropagateError sirve para verificar que cuando ocurre un error en el repositorio,
// el servicio de aplicacion lo propaga correctamente.
func TestTaskService_CreateTask_RepositoryError_ShouldPropagateError(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTaskRepository(ctrl)
	service := application.NewTaskService(mockRepo)
	ctx := context.Background()
	title := "Titulo valido"
	description := "Descripcion valida"

	expectedError := assert.AnError // Crear un error simulado para probar la propagacion

	// configurar mock para simular error del repositorio
	mockRepo.EXPECT().Create(ctx, gomock.Any()).Return(nil, expectedError).Times(1)

	// Act
	result, err := service.CreateTask(ctx, title, description)

	// Assert
	assert.Error(t, err)                // Verificar que se retorna un error
	assert.Nil(t, result)               // Verificar que no se retorna una tarea
	assert.Equal(t, expectedError, err) // Verificar que el error es el esperado

}

// TestTaskService_CreateTask_ValidInputs_ShouldCreateTask verifica que se puede crear una tarea con inputs validos
func TestTaskService_CreateTask_ValidInputs_ShouldCreateTask(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTaskRepository(ctrl)
	service := application.NewTaskService(mockRepo)

	ctx := context.Background()

	// Casos de prueba con diferentes inputs válidos
	testCases := []struct {
		name        string
		title       string
		description string
	}{
		{
			name:        "título y descripción normales",
			title:       "Tarea normal",
			description: "Descripción normal",
		},
		{
			name:        "título y descripción largos",
			title:       "Este es un título muy largo para probar que funciona correctamente",
			description: "Esta es una descripción muy larga para verificar que el sistema maneja correctamente textos extensos",
		},
		{
			name:        "título con caracteres especiales",
			title:       "Tarea #1 - Revisión & Análisis",
			description: "Descripción con símbolos: @, %, $, etc.",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Crear la tarea esperada que devolvería el repositorio
			expectedTask := &domain.Task{
				ID:          1,
				Title:       tc.title,
				Description: tc.description,
				Completed:   false,
				CreatedAt:   time.Now().UTC(),
				UpdatedAt:   time.Now().UTC(),
			}

			// Configurar expectativa del mock para este caso específico
			mockRepo.EXPECT().
				Create(ctx, gomock.Any()).
				Return(expectedTask, nil).
				Times(1)

			// Act
			result, err := service.CreateTask(ctx, tc.title, tc.description)

			// Assert
			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, tc.title, result.Title)
			assert.Equal(t, tc.description, result.Description)
			assert.False(t, result.Completed)
		})
	}
}
