package presentation_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/YerkoTenorio/api-go-hexagonal/modules/task/domain"
	"github.com/YerkoTenorio/api-go-hexagonal/modules/task/presentation"
	"github.com/YerkoTenorio/api-go-hexagonal/modules/task/presentation/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestTaskHandler_GetAllTasks_Success(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockTaskServiceInterface(ctrl)
	handler := presentation.NewTaskHandler(mockService)

	// Datos de prueba

	expectedTasks := []*domain.Task{
		{
			ID:          1,
			Title:       "Tarea 1",
			Description: "Descripcion 1",
			Completed:   false,
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
		},
		{
			ID:          2,
			Title:       "Tarea 2",
			Description: "Descripcion 2",
			Completed:   true,
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
		},
	}

	// Expectativas del mock
	mockService.EXPECT().GetAllTasks(gomock.Any()).Return(expectedTasks, nil).Times(1)

	// Configurar Gin en modo test
	gin.SetMode(gin.TestMode)
	router := gin.New()                       // Crear un router Gin en modo test
	router.GET("/tasks", handler.GetAllTasks) // Registrar la ruta GET /tasks

	// Crear request
	req, _ := http.NewRequest("GET", "/tasks", nil) // Crear una request GET a /tasks
	w := httptest.NewRecorder()                     // Crear un recorder para capturar la respuesta

	// Act
	router.ServeHTTP(w, req) // Servir la request al router

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	//Verificar el JSON de la respuesta
	var response map[string]interface{}              // Crear una variable para almacenar la respuesta JSON
	err := json.Unmarshal(w.Body.Bytes(), &response) // Deserializar el JSON de la respuesta en la variable response
	assert.NoError(t, err)

	// Verificar estructura de la respuesta
	assert.Equal(t, "Tasks retrieved successfully", response["message"]) // Verificar el mensaje de éxito
	assert.Equal(t, float64(2), response["count"])                       // Verificar el conteo de tareas

	// verificar que data existe y es un array
	data, exists := response["data"]
	assert.True(t, exists)
	assert.IsType(t, []interface{}{}, data) // Verificar que data es un array de interfaces

	// Verificar que tenemos 2 tareas en el array
	tasks := data.([]interface{})
	assert.Len(t, tasks, 2) // Verificar que hay 2 tareas en el array

}

// TestTaskHandler_GetAllTasks_ServiceError verifica que se devuelve un error 500 cuando el servicio falla
func TestTaskHandler_GetAllTasks_ServiceError(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockTaskServiceInterface(ctrl)
	handler := presentation.NewTaskHandler(mockService)

	serviceError := errors.New("database connection failed")

	// Expectativas del mock
	mockService.EXPECT().
		GetAllTasks(gomock.Any()).
		Return(nil, serviceError).
		Times(1)

	// Configurar Gin en modo test
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/tasks", handler.GetAllTasks)

	// Crear request
	req, _ := http.NewRequest("GET", "/tasks", nil)
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// Verificar el JSON de respuesta
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Verificar estructura de error
	assert.Equal(t, "Error getting tasks", response["error"])
	assert.Contains(t, response["message"], "database connection failed")
}

func TestTaskHandler_GetAllTasks_EmptyList(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockTaskServiceInterface(ctrl)
	handler := presentation.NewTaskHandler(mockService)

	// Lista vacía
	emptyTasks := []*domain.Task{}

	// Expectativas del mock
	mockService.EXPECT().
		GetAllTasks(gomock.Any()).
		Return(emptyTasks, nil).
		Times(1)

	// Configurar Gin en modo test
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/tasks", handler.GetAllTasks)

	// Crear request
	req, _ := http.NewRequest("GET", "/tasks", nil)
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	// Verificar el JSON de respuesta
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Verificar estructura de respuesta
	assert.Equal(t, "Tasks retrieved successfully", response["message"])
	assert.Equal(t, float64(0), response["count"])

	// Verificar que data existe y es un array vacío
	data, exists := response["data"]
	assert.True(t, exists)
	assert.IsType(t, []interface{}{}, data)

	// Verificar que el array está vacío
	tasks := data.([]interface{})
	assert.Len(t, tasks, 0)
}
