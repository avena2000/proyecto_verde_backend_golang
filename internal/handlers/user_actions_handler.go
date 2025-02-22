package handlers

import (
	"backend_proyecto_verde/internal/models"
	"backend_proyecto_verde/internal/repository/postgres"
	"backend_proyecto_verde/internal/utils"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type UserActionsHandler struct {
	repo         *postgres.UserActionsRepository
	medallasRepo *postgres.MedallasRepository
}

func NewUserActionsHandler(repo *postgres.UserActionsRepository, medallasRepo *postgres.MedallasRepository) *UserActionsHandler {
	return &UserActionsHandler{
		repo:         repo,
		medallasRepo: medallasRepo,
	}
}

func (h *UserActionsHandler) CreateAction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]

	// Parsear el formulario multipart con un límite de 10MB
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Error al procesar el formulario", http.StatusBadRequest)
		return
	}

	// Obtener la imagen del formulario
	file, _, err := r.FormFile("imagen")
	if err != nil && err != http.ErrMissingFile {
		http.Error(w, "Error al obtener la imagen", http.StatusBadRequest)
		return
	}

	// Subir la imagen usando el repositorio
	imageURL, err := h.repo.UploadImage(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
		http.Error(w, "Error al convertir la latitud a float", http.StatusBadRequest)
		return
	}

	longitudFloat, err := strconv.ParseFloat(longitud, 64)
	if err != nil {
		http.Error(w, "Error al convertir la longitud a float", http.StatusBadRequest)
		return
	}

	lugar, ciudad, err := utils.ReverseGeocodeWithCity(latitud, longitud)

	if err != nil {
		http.Error(w, "Error al obtener la ciudad", http.StatusInternalServerError)
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.medallasRepo.AutoAsignMedallas(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(action)
}

func (h *UserActionsHandler) GetUserActions(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]

	actions, err := h.repo.GetUserActions(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(actions)
}

func (h *UserActionsHandler) DeleteAction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Obtener la acción antes de eliminarla para saber el tipo y el userID
	action, err := h.repo.GetActionByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Acción no encontrada", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Eliminar la acción de manera lógica
	if err := h.repo.SoftDeleteAction(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Restar los puntos correspondientes al tipo de acción
	if err := h.repo.DecrementUserPoints(action.UserID, action.TipoAccion); err != nil {
		http.Error(w, "Error al actualizar los puntos del usuario", http.StatusInternalServerError)
		return
	}

	// Decrementar el contador total de acciones del usuario
	if err := h.repo.DecrementTotalActions(action.UserID); err != nil {
		http.Error(w, "Error al actualizar el total de acciones", http.StatusInternalServerError)
		return
	}

	// Verificar y actualizar las medallas del usuario
	if err := h.medallasRepo.VerifyAndUpdateMedallas(action.UserID); err != nil {
		http.Error(w, "Error al verificar las medallas del usuario", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *UserActionsHandler) GetAllActions(w http.ResponseWriter, r *http.Request) {
	actions, err := h.repo.GetAllActions()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(actions)
}
