package presentation_test

import (
	"bytes"
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

// TestTaskHandler_CreateTask_Success verifica la creación exitosa de una tarea
func TestTaskHandler_CreateTask_Success(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockTaskServiceInterface(ctrl)
	handler := presentation.NewTaskHandler(mockService)

	// Datos de entrada
	requestBody := map[string]interface{}{
		"title":       "Nueva Tarea",
		"description": "Descripción de la nueva tarea",
	}

	// Tarea que devuelve el servicio
	expectedTask := &domain.Task{
		ID:          1,
		Title:       "Nueva Tarea",
		Description: "Descripción de la nueva tarea",
		Completed:   false,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	// Expectativas del mock
	mockService.EXPECT().
		CreateTask(gomock.Any(), "Nueva Tarea", "Descripción de la nueva tarea").
		Return(expectedTask, nil).
		Times(1)

	// Configurar Gin en modo test
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/tasks", handler.CreateTask)

	// Crear request con JSON
	jsonBody, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/tasks", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusCreated, w.Code)

	// Verificar el JSON de respuesta
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Verificar estructura de respuesta
	assert.Equal(t, "Task created successfully", response["message"])

	// Verificar que data existe y contiene la tarea
	data, exists := response["data"]
	assert.True(t, exists)
	assert.IsType(t, map[string]interface{}{}, data)

	// Verificar campos de la tarea
	task := data.(map[string]interface{})
	assert.Equal(t, float64(1), task["id"]) // JSON unmarshals numbers as float64
	assert.Equal(t, "Nueva Tarea", task["title"])
	assert.Equal(t, "Descripción de la nueva tarea", task["description"])
	assert.Equal(t, false, task["completed"])
}

// TestTaskHandler_CreateTask_ServiceError verifica el manejo de errores del servicio
func TestTaskHandler_CreateTask_ServiceError(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockTaskServiceInterface(ctrl)
	handler := presentation.NewTaskHandler(mockService)

	// Datos de entrada
	requestBody := map[string]interface{}{
		"title":       "Nueva Tarea",
		"description": "Descripción de la nueva tarea",
	}

	serviceError := errors.New("database connection failed")

	// Expectativas del mock
	mockService.EXPECT().
		CreateTask(gomock.Any(), "Nueva Tarea", "Descripción de la nueva tarea").
		Return(nil, serviceError).
		Times(1)

	// Configurar Gin en modo test
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/tasks", handler.CreateTask)

	// Crear request con JSON
	jsonBody, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/tasks", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
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
	assert.Equal(t, "Error creating task", response["error"])
	assert.Contains(t, response["message"], "database connection failed")
}

// TestTaskHandler_CreateTask_InvalidJSON_MissingTitle verifica validación de título requerido
func TestTaskHandler_CreateTask_InvalidJSON_MissingTitle(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockTaskServiceInterface(ctrl)
	handler := presentation.NewTaskHandler(mockService)

	// JSON sin título
	requestBody := map[string]interface{}{
		"description": "Descripción sin título",
	}

	// NO esperamos llamadas al servicio porque falla la validación antes
	// mockService.EXPECT() - Sin expectativas

	// Configurar Gin en modo test
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/tasks", handler.CreateTask)

	// Crear request con JSON inválido
	jsonBody, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/tasks", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
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
	assert.Equal(t, "Invalid request", response["error"])
	assert.Contains(t, response["message"], "Title")
}

// TestTaskHandler_CreateTask_InvalidJSON_MissingDescription verifica validación de descripción requerida
func TestTaskHandler_CreateTask_InvalidJSON_MissingDescription(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockTaskServiceInterface(ctrl)
	handler := presentation.NewTaskHandler(mockService)

	// JSON sin descripción
	requestBody := map[string]interface{}{
		"title": "Título sin descripción",
	}

	// NO esperamos llamadas al servicio porque falla la validación antes

	// Configurar Gin en modo test
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/tasks", handler.CreateTask)

	// Crear request con JSON inválido
	jsonBody, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/tasks", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
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
	assert.Equal(t, "Invalid request", response["error"])
	assert.Contains(t, response["message"], "Description")
}

// TestTaskHandler_CreateTask_InvalidJSON_EmptyBody verifica manejo de JSON vacío
func TestTaskHandler_CreateTask_InvalidJSON_EmptyBody(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockTaskServiceInterface(ctrl)
	handler := presentation.NewTaskHandler(mockService)

	// NO esperamos llamadas al servicio porque falla la validación antes

	// Configurar Gin en modo test
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/tasks", handler.CreateTask)

	// Crear request con body vacío
	req, _ := http.NewRequest("POST", "/tasks", bytes.NewBuffer([]byte("{}")))
	req.Header.Set("Content-Type", "application/json")
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
	assert.Equal(t, "Invalid request", response["error"])
	// Debe contener error sobre campos requeridos
	message := response["message"].(string)
	assert.True(t,
		(len(message) > 0),
		"El mensaje de error debe contener información sobre campos requeridos")
}

// TestTaskHandler_CreateTask_InvalidJSON_MalformedJSON verifica manejo de JSON malformado
func TestTaskHandler_CreateTask_InvalidJSON_MalformedJSON(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockTaskServiceInterface(ctrl)
	handler := presentation.NewTaskHandler(mockService)

	// NO esperamos llamadas al servicio porque falla la validación antes

	// Configurar Gin en modo test
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/tasks", handler.CreateTask)

	// Crear request con JSON malformado
	req, _ := http.NewRequest("POST", "/tasks", bytes.NewBuffer([]byte(`{"title": "test", "description"`)))
	req.Header.Set("Content-Type", "application/json")
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
	assert.Equal(t, "Invalid request", response["error"])
	assert.NotEmpty(t, response["message"])
}
