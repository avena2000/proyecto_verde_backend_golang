package handlers

import (
	"backend_proyecto_verde/internal/models"
	"backend_proyecto_verde/internal/repository/postgres"
	"backend_proyecto_verde/internal/utils"
	"bytes"
	"database/sql"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gorilla/mux"
)

type UserActionsHandler struct {
	repo         *postgres.UserActionsRepository
	medallasRepo *postgres.MedallasRepository
	awsSession   *session.Session
}

func NewUserActionsHandler(repo *postgres.UserActionsRepository, medallasRepo *postgres.MedallasRepository, awsSession *session.Session) *UserActionsHandler {
	return &UserActionsHandler{
		repo:         repo,
		medallasRepo: medallasRepo,
		awsSession:   awsSession,
	}
}

func (h *UserActionsHandler) CreateAction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]

	// Parsear el formulario multipart con un límite de 10MB
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		utils.RespondWithBadRequest(w, "Error al procesar el formulario", err.Error())
		return
	}

	// Obtener la imagen del formulario
	file, _, err := r.FormFile("imagen")
	if err != nil && err != http.ErrMissingFile {
		utils.RespondWithBadRequest(w, "Error al obtener la imagen", err.Error())
		return
	}

	// Subir la imagen usando el repositorio
	imageURL, err := h.UploadImageToLightsail(file)
	if err != nil {
		utils.RespondWithInternalServerError(w, err.Error(), "Error al subir la imagen")
		return
	}

	tipoAccion := r.FormValue("tipo_accion")
	latitud := r.FormValue("latitud")
	longitud := r.FormValue("longitud")

	// Validar los valores recibidos
	if tipoAccion == "" || latitud == "" || longitud == "" {
		http.Error(w, "Faltan parámetros requeridos", http.StatusBadRequest)
		return
	}

	latitudFloat, err := strconv.ParseFloat(latitud, 64)
	if err != nil {
		utils.RespondWithBadRequest(w, "Latitud inválida", err.Error())
		return
	}

	longitudFloat, err := strconv.ParseFloat(longitud, 64)
	if err != nil {
		utils.RespondWithBadRequest(w, "Longitud inválida", err.Error())
		return
	}

	lugar, ciudad, err := utils.ReverseGeocodeWithCity(latitud, longitud)

	if err != nil {
		utils.RespondWithBadRequest(w, "Error al obtener la ciudad", err.Error())
		return
	}

	action := models.UserAction{
		UserID:         userID,
		TipoAccion:     tipoAccion,
		Latitud:        latitudFloat,
		Longitud:       longitudFloat,
		Foto:           imageURL,
		Ciudad:         ciudad,
		Lugar:          lugar,
		EnColaboracion: false,
		EsParaTorneo:   false,
	}

	err = h.repo.CreateAction(&action)
	if err != nil {
		utils.RespondWithDatabaseError(w, "Error al crear la acción", err.Error())
		return
	}
	// Verificar si el usuario ha ganado medallas
	h.checkMedallas(userID)

	utils.RespondWithCreated(w, action, "Acción creada correctamente")
}

func (h *UserActionsHandler) GetUserActions(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]

	actions, err := h.repo.GetUserActions(userID)
	if err != nil {
		utils.RespondWithDatabaseError(w, "Error al obtener las acciones del usuario", err.Error())
		return
	}

	utils.RespondWithSuccess(w, actions, "Acciones del usuario obtenidas correctamente")
}

func (h *UserActionsHandler) DeleteAction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Obtener la acción antes de eliminarla para saber el tipo y el userID
	action, err := h.repo.GetActionByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.RespondWithNotFound(w, "Acción no encontrada", "No se encontró la acción con el ID proporcionado")
		} else {
			utils.RespondWithDatabaseError(w, "Error al obtener la acción", err.Error())
		}
		return
	}

	// Eliminar la acción
	if err := h.repo.SoftDeleteAction(id); err != nil {
		utils.RespondWithDatabaseError(w, "Error al eliminar la acción", err.Error())
		return
	}

	// Restar los puntos correspondientes al tipo de acción
	if err := h.repo.DecrementUserPoints(action.UserID, action.TipoAccion); err != nil {
		http.Error(w, "Error al actualizar los puntos del usuario", http.StatusInternalServerError)
		return
	}

	// Decrementar el contador total de acciones del usuario
	if err := h.repo.DecrementTotalActions(action.UserID); err != nil {
		// Solo registramos el error pero no interrumpimos la operación
		// ya que la acción ya fue eliminada
		utils.RespondWithSuccess(w, nil, "Acción eliminada correctamente, pero hubo un error al actualizar las estadísticas")
		return
	}

	// Verificar y actualizar las medallas del usuario
	if err := h.medallasRepo.VerifyAndUpdateMedallas(action.UserID); err != nil {
		utils.RespondWithDatabaseError(w, "Error al verificar las medallas del usuario", err.Error())
		return
	}

	utils.RespondWithSuccess(w, nil, "Acción eliminada correctamente")
}

func (h *UserActionsHandler) GetAllActions(w http.ResponseWriter, r *http.Request) {
	actions, err := h.repo.GetAllActions()
	if err != nil {
		utils.RespondWithDatabaseError(w, "Error al obtener todas las acciones", err.Error())
		return
	}

	utils.RespondWithSuccess(w, actions, "Todas las acciones obtenidas correctamente")
}

// Método auxiliar para verificar si el usuario ha ganado medallas
func (h *UserActionsHandler) checkMedallas(userID string) {
	// Esta función se ejecuta en segundo plano y no afecta la respuesta al cliente
	go func() {
		h.medallasRepo.AutoAsignMedallas(userID)
	}()
}

// UploadImageToLightsail sube una imagen a un bucket de AWS usando el SDK de S3
// y devuelve la URL de la imagen subida
func (h *UserActionsHandler) UploadImageToLightsail(file multipart.File) (string, error) {
	if file == nil {
		return "", fmt.Errorf("archivo no proporcionado")
	}

	// Leer el contenido del archivo
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("error al leer el archivo: %v", err)
	}

	// Crear cliente de S3 usando la sesión configurada en el main
	s3Client := s3.New(h.awsSession)

	// Nombre del bucket y clave del objeto
	bucketName := "bucket-e886au" // Reemplaza con el nombre de tu bucket
	objectKey := fmt.Sprintf("images/%s.jpg", time.Now().Format("20060102150405"))

	// Preparar la solicitud para subir el objeto
	input := &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(objectKey),
		Body:        bytes.NewReader(fileBytes),
		ContentType: aws.String("image/jpeg"), // Ajusta según el tipo de archivo
	}

	// Subir el objeto al bucket
	_, err = s3Client.PutObject(input)
	if err != nil {
		return "", fmt.Errorf("error al subir la imagen a S3: %v", err)
	}

	// Construir la URL del objeto
	region := *h.awsSession.Config.Region
	imageURL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucketName, region, objectKey)

	return imageURL, nil
}
