package presentation_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/YerkoTenorio/api-go-hexagonal/modules/task/presentation"
	"github.com/YerkoTenorio/api-go-hexagonal/modules/task/presentation/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

// TestTaskHandler_DeleteTask_Success verifica la eliminación exitosa de una tarea
func TestTaskHandler_DeleteTask_Success(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockTaskServiceInterface(ctrl)
	handler := presentation.NewTaskHandler(mockService)

	// Expectativas del mock
	mockService.EXPECT().
		DeleteTask(gomock.Any(), 1).
		Return(nil).
		Times(1)

	// Configurar Gin en modo test
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.DELETE("/tasks/:id", handler.DeleteTask)

	// Crear request
	req, _ := http.NewRequest("DELETE", "/tasks/1", nil)
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
	assert.Equal(t, "Task deleted successfully", response["message"])

	// Verificar que NO hay campo "data" (a diferencia de otros endpoints)
	_, hasData := response["data"]
	assert.False(t, hasData, "Delete endpoint should not return data field")
}

// TestTaskHandler_DeleteTask_ServiceError verifica el manejo de errores del servicio (tarea no encontrada)
func TestTaskHandler_DeleteTask_ServiceError(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockTaskServiceInterface(ctrl)
	handler := presentation.NewTaskHandler(mockService)

	serviceError := errors.New("task not found")

	// Expectativas del mock
	mockService.EXPECT().
		DeleteTask(gomock.Any(), 999).
		Return(serviceError).
		Times(1)

	// Configurar Gin en modo test
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.DELETE("/tasks/:id", handler.DeleteTask)

	// Crear request con ID que no existe
	req, _ := http.NewRequest("DELETE", "/tasks/999", nil)
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
	assert.Equal(t, "Error deleting task", response["error"])
	assert.Contains(t, response["message"], "task not found")
}

// TestTaskHandler_DeleteTask_InvalidID_NonNumeric verifica validación de ID no numérico
func TestTaskHandler_DeleteTask_InvalidID_NonNumeric(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockTaskServiceInterface(ctrl)
	handler := presentation.NewTaskHandler(mockService)

	// NO esperamos llamadas al servicio porque falla la validación antes

	// Configurar Gin en modo test
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.DELETE("/tasks/:id", handler.DeleteTask)

	// Crear request con ID no numérico
	req, _ := http.NewRequest("DELETE", "/tasks/abc", nil)
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

// TestTaskHandler_DeleteTask_InvalidID_Negative verifica validación de ID negativo
func TestTaskHandler_DeleteTask_InvalidID_Negative(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockTaskServiceInterface(ctrl)
	handler := presentation.NewTaskHandler(mockService)

	// NO esperamos llamadas al servicio porque falla la validación antes

	// Configurar Gin en modo test
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.DELETE("/tasks/:id", handler.DeleteTask)

	// Crear request con ID negativo
	req, _ := http.NewRequest("DELETE", "/tasks/-1", nil)
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

// TestTaskHandler_DeleteTask_InvalidID_Zero verifica validación de ID cero
func TestTaskHandler_DeleteTask_InvalidID_Zero(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockTaskServiceInterface(ctrl)
	handler := presentation.NewTaskHandler(mockService)

	// Expectativas del mock - ID 0 es técnicamente válido para ParseUint pero podría no existir
	mockService.EXPECT().
		DeleteTask(gomock.Any(), 0).
		Return(errors.New("task not found")).
		Times(1)

	// Configurar Gin en modo test
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.DELETE("/tasks/:id", handler.DeleteTask)

	// Crear request con ID cero
	req, _ := http.NewRequest("DELETE", "/tasks/0", nil)
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
	assert.Equal(t, "Error deleting task", response["error"])
	assert.Contains(t, response["message"], "task not found")
}

// TestTaskHandler_DeleteTask_InvalidID_Float verifica validación de ID decimal
func TestTaskHandler_DeleteTask_InvalidID_Float(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockTaskServiceInterface(ctrl)
	handler := presentation.NewTaskHandler(mockService)

	// NO esperamos llamadas al servicio porque falla la validación antes

	// Configurar Gin en modo test
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.DELETE("/tasks/:id", handler.DeleteTask)

	// Crear request con ID decimal
	req, _ := http.NewRequest("DELETE", "/tasks/1.5", nil)
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

// TestTaskHandler_DeleteTask_AlreadyDeleted verifica el comportamiento al intentar eliminar una tarea ya eliminada
func TestTaskHandler_DeleteTask_AlreadyDeleted(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockTaskServiceInterface(ctrl)
	handler := presentation.NewTaskHandler(mockService)

	serviceError := errors.New("task not found")

	// Expectativas del mock
	mockService.EXPECT().
		DeleteTask(gomock.Any(), 1).
		Return(serviceError).
		Times(1)

	// Configurar Gin en modo test
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.DELETE("/tasks/:id", handler.DeleteTask)

	// Crear request
	req, _ := http.NewRequest("DELETE", "/tasks/1", nil)
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
	assert.Equal(t, "Error deleting task", response["error"])
	assert.Contains(t, response["message"], "task not found")
}
