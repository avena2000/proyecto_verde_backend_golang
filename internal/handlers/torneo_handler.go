package handlers

import (
	"backend_proyecto_verde/internal/models"
	"backend_proyecto_verde/internal/repository/postgres"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type TorneoHandler struct {
	repo *postgres.TorneoRepository
}

func NewTorneoHandler(repo *postgres.TorneoRepository) *TorneoHandler {
	return &TorneoHandler{repo: repo}
}

func (h *TorneoHandler) CreateTorneo(w http.ResponseWriter, r *http.Request) {
	var torneo models.Torneo
	if err := json.NewDecoder(r.Body).Decode(&torneo); err != nil {
		http.Error(w, "Error al decodificar el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}

	if err := h.repo.CreateTorneo(&torneo); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(torneo)
}

func (h *TorneoHandler) GetTorneo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	torneo, err := h.repo.GetTorneoByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(torneo)
}

func (h *TorneoHandler) UpdateTorneo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var torneo models.Torneo
	if err := json.NewDecoder(r.Body).Decode(&torneo); err != nil {
		http.Error(w, "Error al decodificar el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}

	torneo.ID = id
	if err := h.repo.UpdateTorneo(&torneo); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(torneo)
}

func (h *TorneoHandler) ListTorneos(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	if limit <= 0 {
		limit = 10
	}

	torneos, err := h.repo.ListTorneos(limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(torneos)
}

func (h *TorneoHandler) GetTorneoStats(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	//id := vars["id"]

	// Implementar lógica para obtener estadísticas
	w.Header().Set("Content-Type", "application/json")
	// TODO: Implementar obtención de estadísticas
} 