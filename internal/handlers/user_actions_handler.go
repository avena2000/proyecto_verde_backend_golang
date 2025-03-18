package handlers

import (
	"backend_proyecto_verde/internal/models"
	"backend_proyecto_verde/internal/repository/postgres"
	"backend_proyecto_verde/internal/utils"
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"

	"git.sr.ht/~jamesponddotco/bunnystorage-go"
	"github.com/gorilla/mux"
)

type UserActionsHandler struct {
	repo         *postgres.UserActionsRepository
	medallasRepo *postgres.MedallasRepository
	bunnyClient  *bunnystorage.Client
	storageZone  string
}

func NewUserActionsHandler(repo *postgres.UserActionsRepository, medallasRepo *postgres.MedallasRepository, bunnyClient *bunnystorage.Client, storageZone string) *UserActionsHandler {
	return &UserActionsHandler{
		repo:         repo,
		medallasRepo: medallasRepo,
		bunnyClient:  bunnyClient,
		storageZone:  storageZone,
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
	imageURL, err := h.UploadImageToBunnyStorage(file)
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

// UploadImageToBunnyStorage sube una imagen a BunnyStorage usando el SDK de bunnystorage-go
// y devuelve la URL de la imagen subida
func (h *UserActionsHandler) UploadImageToBunnyStorage(file multipart.File) (string, error) {
	if file == nil {
		return "", fmt.Errorf("archivo no proporcionado")
	}

	// Leer el contenido del archivo
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("error al leer el archivo: %v", err)
	}

	// Verificar que el archivo tenga contenido
	if len(fileBytes) == 0 {
		return "", fmt.Errorf("el archivo está vacío")
	}

	// Imprimir información para depuración
	fmt.Printf("Tamaño del archivo a subir: %d bytes\n", len(fileBytes))
	fmt.Printf("Cliente BunnyStorage inicializado: %v\n", h.bunnyClient != nil)
	fmt.Printf("Zona de almacenamiento: %s\n", h.storageZone)

	// Generar la ruta del archivo
	fileName := fmt.Sprintf("%s.jpg", time.Now().Format("20060102150405"))
	fmt.Printf("Nombre del archivo a subir: %s\n", fileName)

	// Crear el contexto para la operación
	ctx := context.Background()

	// Subir el archivo a BunnyStorage con manejo de errores detallado
	fmt.Println("Iniciando carga a BunnyStorage...")
	upload, err := h.bunnyClient.Upload(ctx, "/images", fileName, "", bytes.NewReader(fileBytes))
	if err != nil {
		fmt.Printf("Error al subir la imagen a BunnyStorage: %v\n", err)
		return "", fmt.Errorf("error al subir la imagen a BunnyStorage: %v", err)
	}

	// Imprimir información del resultado
	fmt.Printf("Archivo subido con éxito. Detalles: %+v\n", upload)

	// Construir la URL del objeto
	imageURL := fmt.Sprintf("https://%s.b-cdn.net/images/%s", h.storageZone, fileName)
	fmt.Printf("URL generada: %s\n", imageURL)

	return imageURL, nil
}
