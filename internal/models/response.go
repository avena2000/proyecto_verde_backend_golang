package models

// APIResponse es la estructura estandarizada para todas las respuestas de la API
type APIResponse struct {
	Code    string      `json:"code"`    // Código que indica el resultado de la operación (000 para éxito)
	Message string      `json:"message"` // Mensaje personalizado, puede ser null
	Description string    `json:"description"` // Descripción del error, puede ser null
	Data    interface{} `json:"data"`    // Datos de respuesta, puede ser cualquier modelo o null
}

// Códigos de respuesta comunes
const (
	// Códigos de éxito
	CodeSuccess = "000" // Operación exitosa

	// Códigos de error de cliente (4xx)
	CodeBadRequest       = "400" // Solicitud incorrecta
	CodeUnauthorized     = "401" // No autorizado
	CodeForbidden        = "403" // Prohibido
	CodeNotFound         = "404" // Recurso no encontrado
	CodeMethodNotAllowed = "405" // Método no permitido
	CodeConflict         = "409" // Conflicto con el estado actual del recurso
	CodeValidationError  = "422" // Error de validación

	// Códigos de error de servidor (5xx)
	CodeInternalServerError = "500" // Error interno del servidor
	CodeServiceUnavailable  = "503" // Servicio no disponible
	CodeDatabaseError       = "510" // Error de base de datos

	// Códigos personalizados para la aplicación
	CodeUserAlreadyExists  = "601" // El usuario ya existe
	CodeInvalidCredentials = "602" // Credenciales inválidas
	CodeResourceNotCreated = "603" // No se pudo crear el recurso
	CodeResourceNotUpdated = "604" // No se pudo actualizar el recurso
	CodeResourceNotDeleted = "605" // No se pudo eliminar el recurso
	CodeImageUploadError   = "606" // Error al subir la imagen
	CodeTournamentError    = "607" // Error relacionado con torneos
	CodeFriendshipError    = "608" // Error relacionado con amistades
	CodeMedalError         = "609" // Error relacionado con medallas
)

// NewSuccessResponse crea una nueva respuesta de éxito
func NewSuccessResponse(data interface{}, message string) APIResponse {
	return APIResponse{
		Code:    CodeSuccess,
		Message: message,
		Data:    data,
	}
}

// NewErrorResponse crea una nueva respuesta de error
func NewErrorResponse(code string, message string, description string) APIResponse {
	return APIResponse{
		Code:    code,
		Message: message,
		Description: description,
		Data:    nil,
	}
}
