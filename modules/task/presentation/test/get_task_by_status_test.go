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

// TestTaskHandler_GetTaskByStatus_Success_Completed verifica la obtención exitosa de tareas completadas
func TestTaskHandler_GetTaskByStatus_Success_Completed(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockTaskServiceInterface(ctrl)
	handler := presentation.NewTaskHandler(mockService)

	// Tareas completadas que devuelve el servicio
	expectedTasks := []*domain.Task{
		{
			ID:          1,
			Title:       "Tarea Completada 1",
			Description: "Descripción 1",
			Completed:   true,
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
		},
		{
			ID:          2,
			Title:       "Tarea Completada 2",
			Description: "Descripción 2",
			Completed:   true,
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
		},
	}

	// Expectativas del mock
	mockService.EXPECT().
		GetTasksByStatus(gomock.Any(), true).
		Return(expectedTasks, nil).
		Times(1)

	// Configurar Gin en modo test
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/tasks/status", handler.GetTaskByStatus)

	// Crear request con query parameter completed=true
	req, _ := http.NewRequest("GET", "/tasks/status?completed=true", nil)
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

	// Verificar que data existe y es un array
	data, exists := response["data"]
	assert.True(t, exists)
	assert.IsType(t, []interface{}{}, data)

	// Verificar que tenemos 2 tareas
	tasks := data.([]interface{})
	assert.Len(t, tasks, 2)

	// Verificar que todas las tareas están completadas
	for _, taskInterface := range tasks {
		task := taskInterface.(map[string]interface{})
		assert.Equal(t, true, task["completed"])
	}
}

// TestTaskHandler_GetTaskByStatus_Success_NotCompleted verifica la obtención exitosa de tareas no completadas
func TestTaskHandler_GetTaskByStatus_Success_NotCompleted(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockTaskServiceInterface(ctrl)
	handler := presentation.NewTaskHandler(mockService)

	// Tareas no completadas que devuelve el servicio
	expectedTasks := []*domain.Task{
		{
			ID:          3,
			Title:       "Tarea Pendiente 1",
			Description: "Descripción 3",
			Completed:   false,
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
		},
	}

	// Expectativas del mock
	mockService.EXPECT().
		GetTasksByStatus(gomock.Any(), false).
		Return(expectedTasks, nil).
		Times(1)

	// Configurar Gin en modo test
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/tasks/status", handler.GetTaskByStatus)

	// Crear request con query parameter completed=false
	req, _ := http.NewRequest("GET", "/tasks/status?completed=false", nil)
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

	// Verificar que data existe y es un array
	data, exists := response["data"]
	assert.True(t, exists)
	assert.IsType(t, []interface{}{}, data)

	// Verificar que tenemos 1 tarea
	tasks := data.([]interface{})
	assert.Len(t, tasks, 1)

	// Verificar que la tarea no está completada
	task := tasks[0].(map[string]interface{})
	assert.Equal(t, false, task["completed"])
}

// TestTaskHandler_GetTaskByStatus_Success_EmptyList verifica el comportamiento con lista vacía
func TestTaskHandler_GetTaskByStatus_Success_EmptyList(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockTaskServiceInterface(ctrl)
	handler := presentation.NewTaskHandler(mockService)

	// Lista vacía
	emptyTasks := []*domain.Task{}

	// Expectativas del mock
	mockService.EXPECT().
		GetTasksByStatus(gomock.Any(), true).
		Return(emptyTasks, nil).
		Times(1)

	// Configurar Gin en modo test
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/tasks/status", handler.GetTaskByStatus)

	// Crear request
	req, _ := http.NewRequest("GET", "/tasks/status?completed=true", nil)
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

	// Verificar que data existe y es un array vacío
	data, exists := response["data"]
	assert.True(t, exists)
	assert.IsType(t, []interface{}{}, data)

	// Verificar que el array está vacío
	tasks := data.([]interface{})
	assert.Len(t, tasks, 0)
}

// TestTaskHandler_GetTaskByStatus_ServiceError verifica el manejo de errores del servicio
func TestTaskHandler_GetTaskByStatus_ServiceError(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockTaskServiceInterface(ctrl)
	handler := presentation.NewTaskHandler(mockService)

	serviceError := errors.New("database connection failed")

	// Expectativas del mock
	mockService.EXPECT().
		GetTasksByStatus(gomock.Any(), true).
		Return(nil, serviceError).
		Times(1)

	// Configurar Gin en modo test
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/tasks/status", handler.GetTaskByStatus)

	// Crear request
	req, _ := http.NewRequest("GET", "/tasks/status?completed=true", nil)
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code) // Nota: 500, no 400

	// Verificar el JSON de respuesta
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Verificar estructura de error
	assert.Equal(t, "Error getting tasks", response["error"])
	assert.Contains(t, response["message"], "database connection failed")
}

// TestTaskHandler_GetTaskByStatus_InvalidParameter_NonBoolean verifica validación de parámetro no booleano
func TestTaskHandler_GetTaskByStatus_InvalidParameter_NonBoolean(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockTaskServiceInterface(ctrl)
	handler := presentation.NewTaskHandler(mockService)

	// NO esperamos llamadas al servicio porque falla la validación antes

	// Configurar Gin en modo test
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/tasks/status", handler.GetTaskByStatus)

	// Crear request con parámetro no booleano
	req, _ := http.NewRequest("GET", "/tasks/status?completed=maybe", nil)
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Verificar el JSON de respuesta
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Verificar estructura de error
	assert.Equal(t, "Invalid completed parameter", response["error"])
	assert.Equal(t, "completed must be a boolean value", response["message"])
}

// TestTaskHandler_GetTaskByStatus_InvalidParameter_Empty verifica validación de parámetro vacío
func TestTaskHandler_GetTaskByStatus_InvalidParameter_Empty(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockTaskServiceInterface(ctrl)
	handler := presentation.NewTaskHandler(mockService)

	// NO esperamos llamadas al servicio porque falla la validación antes

	// Configurar Gin en modo test
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/tasks/status", handler.GetTaskByStatus)

	// Crear request sin parámetro completed
	req, _ := http.NewRequest("GET", "/tasks/status", nil)
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Verificar el JSON de respuesta
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Verificar estructura de error
	assert.Equal(t, "Invalid completed parameter", response["error"])
	assert.Equal(t, "completed must be a boolean value", response["message"])
}

// TestTaskHandler_GetTaskByStatus_ValidParameter_Numeric verifica que "1" es válido como true
func TestTaskHandler_GetTaskByStatus_ValidParameter_Numeric(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockTaskServiceInterface(ctrl)
	handler := presentation.NewTaskHandler(mockService)

	// Lista vacía para simplificar
	emptyTasks := []*domain.Task{}

	// SÍ esperamos llamadas al servicio porque "1" es válido como true
	mockService.EXPECT().
		GetTasksByStatus(gomock.Any(), true).
		Return(emptyTasks, nil).
		Times(1)

	// Configurar Gin en modo test
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/tasks/status", handler.GetTaskByStatus)

	// Crear request con parámetro numérico "1" (válido como true)
	req, _ := http.NewRequest("GET", "/tasks/status?completed=1", nil)
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	// Verificar el JSON de respuesta
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Verificar estructura de respuesta exitosa
	assert.Equal(t, "Tasks retrieved successfully", response["message"])
	assert.NotNil(t, response["data"])
}
