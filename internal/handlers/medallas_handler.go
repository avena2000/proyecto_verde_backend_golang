package handlers

import (
	"backend_proyecto_verde/internal/models"
	"backend_proyecto_verde/internal/repository/postgres"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

type MedallasHandler struct {
	repo *postgres.MedallasRepository
}

func NewMedallasHandler(repo *postgres.MedallasRepository) *MedallasHandler {
	return &MedallasHandler{repo: repo}
}

func (h *MedallasHandler) CreateMedalla(w http.ResponseWriter, r *http.Request) {
	var medalla models.Medalla
	if err := json.NewDecoder(r.Body).Decode(&medalla); err != nil {
		http.Error(w, "Error al decodificar el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}

	if err := h.repo.CreateMedalla(&medalla); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(medalla)
}

func (h *MedallasHandler) GetMedallas(w http.ResponseWriter, r *http.Request) {
	medallas, err := h.repo.GetMedallas()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(medallas)
}

func (h *MedallasHandler) AsignarMedalla(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]
	medallaID := vars["medalla_id"]

	if err := h.repo.AsignarMedalla(userID, medallaID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *MedallasHandler) GetMedallasUsuario(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]

	medallas, err := h.repo.GetMedallasUsuario(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(medallas)
} 