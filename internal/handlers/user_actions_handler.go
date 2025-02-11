package handlers

import (
	"backend_proyecto_verde/internal/models"
	"backend_proyecto_verde/internal/repository/postgres"
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

	action := models.UserAction{
		UserID:         userID,
		TipoAccion:     tipoAccion,
		Latitud:        latitudFloat,
		Longitud:       longitudFloat,
		Foto:           imageURL,
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

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	if limit <= 0 {
		limit = 10
	}

	actions, err := h.repo.GetUserActions(userID, limit, offset)
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

	if err := h.repo.SoftDeleteAction(id); err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Acción no encontrada", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
