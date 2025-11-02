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

// TestTaskHandler_UpdateTask_Success verifica la actualización exitosa de una tarea
func TestTaskHandler_UpdateTask_Success(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockTaskServiceInterface(ctrl)
	handler := presentation.NewTaskHandler(mockService)

	// Datos de entrada
	requestBody := map[string]interface{}{
		"title":       "Tarea Actualizada",
		"description": "Descripción actualizada",
		"completed":   true,
	}

	// Tarea que devuelve el servicio
	expectedTask := &domain.Task{
		ID:          1,
		Title:       "Tarea Actualizada",
		Description: "Descripción actualizada",
		Completed:   true,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	// Expectativas del mock
	completed := true
	mockService.EXPECT().
		UpdateTask(gomock.Any(), 1, "Tarea Actualizada", "Descripción actualizada", &completed).
		Return(expectedTask, nil).
		Times(1)

	// Configurar Gin en modo test
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.PUT("/tasks/:id", handler.UpdateTask)

	// Crear request con JSON
	jsonBody, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("PUT", "/tasks/1", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
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
	assert.Equal(t, "Task updated successfully", response["message"])

	// Verificar que data existe y contiene la tarea
	data, exists := response["data"]
	assert.True(t, exists)
	assert.IsType(t, map[string]interface{}{}, data)

	// Verificar campos de la tarea
	task := data.(map[string]interface{})
	assert.Equal(t, float64(1), task["id"])
	assert.Equal(t, "Tarea Actualizada", task["title"])
	assert.Equal(t, "Descripción actualizada", task["description"])
	assert.Equal(t, true, task["completed"])
}

// TestTaskHandler_UpdateTask_PartialUpdate verifica actualización parcial (solo algunos campos)
func TestTaskHandler_UpdateTask_PartialUpdate(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockTaskServiceInterface(ctrl)
	handler := presentation.NewTaskHandler(mockService)

	// Solo actualizamos el título
	requestBody := map[string]interface{}{
		"title": "Solo Título Actualizado",
	}

	// Tarea que devuelve el servicio
	expectedTask := &domain.Task{
		ID:          1,
		Title:       "Solo Título Actualizado",
		Description: "Descripción original", // No cambió
		Completed:   false,                  // No cambió
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	// Expectativas del mock - completed será false por defecto
	completed := false
	mockService.EXPECT().
		UpdateTask(gomock.Any(), 1, "Solo Título Actualizado", "", &completed).
		Return(expectedTask, nil).
		Times(1)

	// Configurar Gin en modo test
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.PUT("/tasks/:id", handler.UpdateTask)

	// Crear request con JSON parcial
	jsonBody, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("PUT", "/tasks/1", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
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
	assert.Equal(t, "Task updated successfully", response["message"])

	// Verificar que data existe
	data, exists := response["data"]
	assert.True(t, exists)

	// Verificar que el título se actualizó
	task := data.(map[string]interface{})
	assert.Equal(t, "Solo Título Actualizado", task["title"])
}

// TestTaskHandler_UpdateTask_ServiceError verifica el manejo de errores del servicio
func TestTaskHandler_UpdateTask_ServiceError(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockTaskServiceInterface(ctrl)
	handler := presentation.NewTaskHandler(mockService)

	// Datos de entrada
	requestBody := map[string]interface{}{
		"title": "Tarea que fallará",
	}

	serviceError := errors.New("task not found")

	// Expectativas del mock
	completed := false
	mockService.EXPECT().
		UpdateTask(gomock.Any(), 999, "Tarea que fallará", "", &completed).
		Return(nil, serviceError).
		Times(1)

	// Configurar Gin en modo test
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.PUT("/tasks/:id", handler.UpdateTask)

	// Crear request con JSON
	jsonBody, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("PUT", "/tasks/999", bytes.NewBuffer(jsonBody))
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
	assert.Equal(t, "Error updating task", response["error"])
	assert.Contains(t, response["message"], "task not found")
}

// TestTaskHandler_UpdateTask_InvalidID_NonNumeric verifica validación de ID no numérico
func TestTaskHandler_UpdateTask_InvalidID_NonNumeric(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockTaskServiceInterface(ctrl)
	handler := presentation.NewTaskHandler(mockService)

	// Datos de entrada
	requestBody := map[string]interface{}{
		"title": "No importa",
	}

	// NO esperamos llamadas al servicio porque falla la validación antes

	// Configurar Gin en modo test
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.PUT("/tasks/:id", handler.UpdateTask)

	// Crear request con ID no numérico
	jsonBody, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("PUT", "/tasks/abc", bytes.NewBuffer(jsonBody))
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
	assert.Equal(t, "Invalid ID", response["error"])
	assert.Equal(t, "ID must be a positive integer", response["message"])
}

// TestTaskHandler_UpdateTask_InvalidJSON_MalformedJSON verifica manejo de JSON malformado
func TestTaskHandler_UpdateTask_InvalidJSON_MalformedJSON(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockTaskServiceInterface(ctrl)
	handler := presentation.NewTaskHandler(mockService)

	// NO esperamos llamadas al servicio porque falla la validación antes

	// Configurar Gin en modo test
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.PUT("/tasks/:id", handler.UpdateTask)

	// Crear request con JSON malformado
	req, _ := http.NewRequest("PUT", "/tasks/1", bytes.NewBuffer([]byte(`{"title": "test", "description"`)))
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

// TestTaskHandler_UpdateTask_EmptyJSON verifica manejo de JSON vacío (debería ser válido)
func TestTaskHandler_UpdateTask_EmptyJSON(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockTaskServiceInterface(ctrl)
	handler := presentation.NewTaskHandler(mockService)

	// Tarea que devuelve el servicio (sin cambios)
	expectedTask := &domain.Task{
		ID:          1,
		Title:       "Título Original",
		Description: "Descripción Original",
		Completed:   false,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	// Expectativas del mock - todos los campos vacíos
	completed := false
	mockService.EXPECT().
		UpdateTask(gomock.Any(), 1, "", "", &completed).
		Return(expectedTask, nil).
		Times(1)

	// Configurar Gin en modo test
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.PUT("/tasks/:id", handler.UpdateTask)

	// Crear request con JSON vacío
	req, _ := http.NewRequest("PUT", "/tasks/1", bytes.NewBuffer([]byte("{}")))
	req.Header.Set("Content-Type", "application/json")
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
	assert.Equal(t, "Task updated successfully", response["message"])
	assert.NotNil(t, response["data"])
}
