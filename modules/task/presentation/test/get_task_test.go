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

// TestTaskHandler_GetTask_Success verifica la obtención exitosa de una tarea por ID
func TestTaskHandler_GetTask_Success(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockTaskServiceInterface(ctrl)
	handler := presentation.NewTaskHandler(mockService)

	// Tarea que devuelve el servicio
	expectedTask := &domain.Task{
		ID:          1,
		Title:       "Tarea de Prueba",
		Description: "Descripción de la tarea de prueba",
		Completed:   false,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	// Expectativas del mock
	mockService.EXPECT().
		GetTaskByID(gomock.Any(), 1).
		Return(expectedTask, nil).
		Times(1)

	// Configurar Gin en modo test
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/tasks/:id", handler.GetTask)

	// Crear request
	req, _ := http.NewRequest("GET", "/tasks/1", nil)
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
	assert.Equal(t, "Task retrieved successfully", response["message"])

	// Verificar que data existe y contiene la tarea
	data, exists := response["data"]
	assert.True(t, exists)
	assert.IsType(t, map[string]interface{}{}, data)

	// Verificar campos de la tarea
	task := data.(map[string]interface{})
	assert.Equal(t, float64(1), task["id"]) // JSON unmarshals numbers as float64
	assert.Equal(t, "Tarea de Prueba", task["title"])
	assert.Equal(t, "Descripción de la tarea de prueba", task["description"])
	assert.Equal(t, false, task["completed"])
}

// TestTaskHandler_GetTask_ServiceError verifica el manejo de errores del servicio (tarea no encontrada)
func TestTaskHandler_GetTask_ServiceError(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockTaskServiceInterface(ctrl)
	handler := presentation.NewTaskHandler(mockService)

	serviceError := errors.New("task not found")

	// Expectativas del mock
	mockService.EXPECT().
		GetTaskByID(gomock.Any(), 999).
		Return(nil, serviceError).
		Times(1)

	// Configurar Gin en modo test
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/tasks/:id", handler.GetTask)

	// Crear request con ID que no existe
	req, _ := http.NewRequest("GET", "/tasks/999", nil)
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
	assert.Equal(t, "Error getting task", response["error"])
	assert.Contains(t, response["message"], "task not found")
}

// TestTaskHandler_GetTask_InvalidID_NonNumeric verifica validación de ID no numérico
func TestTaskHandler_GetTask_InvalidID_NonNumeric(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockTaskServiceInterface(ctrl)
	handler := presentation.NewTaskHandler(mockService)

	// NO esperamos llamadas al servicio porque falla la validación antes

	// Configurar Gin en modo test
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/tasks/:id", handler.GetTask)

	// Crear request con ID no numérico
	req, _ := http.NewRequest("GET", "/tasks/abc", nil)
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
	assert.Equal(t, "Invalid ID", response["error"])
	assert.Equal(t, "ID must be a positive integer", response["message"])
}

// TestTaskHandler_GetTask_InvalidID_Negative verifica validación de ID negativo
func TestTaskHandler_GetTask_InvalidID_Negative(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockTaskServiceInterface(ctrl)
	handler := presentation.NewTaskHandler(mockService)

	// NO esperamos llamadas al servicio porque falla la validación antes

	// Configurar Gin en modo test
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/tasks/:id", handler.GetTask)

	// Crear request con ID negativo
	req, _ := http.NewRequest("GET", "/tasks/-1", nil)
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
	assert.Equal(t, "Invalid ID", response["error"])
	assert.Equal(t, "ID must be a positive integer", response["message"])
}

// TestTaskHandler_GetTask_InvalidID_Zero verifica validación de ID cero
func TestTaskHandler_GetTask_InvalidID_Zero(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockTaskServiceInterface(ctrl)
	handler := presentation.NewTaskHandler(mockService)

	// Expectativas del mock - ID 0 es técnicamente válido para ParseUint pero podría no existir
	mockService.EXPECT().
		GetTaskByID(gomock.Any(), 0).
		Return(nil, errors.New("task not found")).
		Times(1)

	// Configurar Gin en modo test
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/tasks/:id", handler.GetTask)

	// Crear request con ID cero
	req, _ := http.NewRequest("GET", "/tasks/0", nil)
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
	assert.Equal(t, "Error getting task", response["error"])
	assert.Contains(t, response["message"], "task not found")
}

// TestTaskHandler_GetTask_InvalidID_Float verifica validación de ID decimal
func TestTaskHandler_GetTask_InvalidID_Float(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockTaskServiceInterface(ctrl)
	handler := presentation.NewTaskHandler(mockService)

	// NO esperamos llamadas al servicio porque falla la validación antes

	// Configurar Gin en modo test
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/tasks/:id", handler.GetTask)

	// Crear request con ID decimal
	req, _ := http.NewRequest("GET", "/tasks/1.5", nil)
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
	assert.Equal(t, "Invalid ID", response["error"])
	assert.Equal(t, "ID must be a positive integer", response["message"])
}
