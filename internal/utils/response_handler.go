package utils

import (
	"backend_proyecto_verde/internal/models"
	"encoding/json"
	"log"
	"net/http"
)

// RespondWithJSON envía una respuesta JSON con el código de estado HTTP proporcionado
func RespondWithJSON(w http.ResponseWriter, statusCode int, response models.APIResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error al codificar la respuesta JSON: %v", err)
	}
}

// RespondWithSuccess envía una respuesta de éxito con los datos proporcionados
func RespondWithSuccess(w http.ResponseWriter, data interface{}, message string) {
	response := models.NewSuccessResponse(data, message)
	RespondWithJSON(w, http.StatusOK, response)
}

// RespondWithCreated envía una respuesta de recurso creado con los datos proporcionados
func RespondWithCreated(w http.ResponseWriter, data interface{}, message string) {
	response := models.NewSuccessResponse(data, message)
	RespondWithJSON(w, http.StatusCreated, response)
}

// RespondWithError envía una respuesta de error con el código y mensaje proporcionados
func RespondWithError(w http.ResponseWriter, statusCode int, errorCode string, message string, description string) {
	response := models.NewErrorResponse(errorCode, message, description)
	RespondWithJSON(w, statusCode, response)
}

// RespondWithBadRequest envía una respuesta de solicitud incorrecta
func RespondWithBadRequest(w http.ResponseWriter, message string, description string) {
	RespondWithError(w, http.StatusBadRequest, models.CodeBadRequest, message, description)
}

// RespondWithUnauthorized envía una respuesta de no autorizado
func RespondWithUnauthorized(w http.ResponseWriter, message string, description string) {
	RespondWithError(w, http.StatusUnauthorized, models.CodeUnauthorized, message, description)
}

// RespondWithForbidden envía una respuesta de prohibido
func RespondWithForbidden(w http.ResponseWriter, message string, description string) {
	RespondWithError(w, http.StatusForbidden, models.CodeForbidden, message, description)
}

// RespondWithNotFound envía una respuesta de recurso no encontrado
func RespondWithNotFound(w http.ResponseWriter, message string, description string) {
	RespondWithError(w, http.StatusNotFound, models.CodeNotFound, message, description)
}

// RespondWithConflict envía una respuesta de conflicto
func RespondWithConflict(w http.ResponseWriter, message string, description string) {
	RespondWithError(w, http.StatusConflict, models.CodeConflict, message, description)
}

// RespondWithInternalServerError envía una respuesta de error interno del servidor
func RespondWithInternalServerError(w http.ResponseWriter, message string, description string) {
	RespondWithError(w, http.StatusInternalServerError, models.CodeInternalServerError, message, description)
}

// RespondWithValidationError envía una respuesta de error de validación
func RespondWithValidationError(w http.ResponseWriter, message string, description string) {
	RespondWithError(w, http.StatusUnprocessableEntity, models.CodeValidationError, message, description)
}

// RespondWithDatabaseError envía una respuesta de error de base de datos
func RespondWithDatabaseError(w http.ResponseWriter, message string, description string) {
	RespondWithError(w, http.StatusConflict, models.CodeDatabaseError, message, description)
}

// RespondWithUserAlreadyExists envía una respuesta de usuario ya existente
func RespondWithUserAlreadyExists(w http.ResponseWriter, message string, description string) {
	RespondWithError(w, http.StatusConflict, models.CodeUserAlreadyExists, message, description)
}

// RespondWithInvalidCredentials envía una respuesta de credenciales inválidas
func RespondWithInvalidCredentials(w http.ResponseWriter, message string, description string) {
	RespondWithError(w, http.StatusUnauthorized, models.CodeInvalidCredentials, message, description)
}
