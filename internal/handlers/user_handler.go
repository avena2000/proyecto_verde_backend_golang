package handlers

import (
	"backend_proyecto_verde/internal/models"
	"backend_proyecto_verde/internal/repository/postgres"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type UserHandler struct {
	repo *postgres.UserRepository
}

func NewUserHandler(repo *postgres.UserRepository) *UserHandler {
	return &UserHandler{repo: repo}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var createUser models.CreateUserAccess
	if err := json.NewDecoder(r.Body).Decode(&createUser); err != nil {
		http.Error(w, "Error al decodificar el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}

	existingUser, err := h.repo.GetUserByUsername(createUser.Username)
	if err == nil && existingUser != nil {
		http.Error(w, "El usuario ya existe", http.StatusConflict)
		return
	}

	user, err := h.repo.CreateUser(&createUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	stats := models.UserStats{
		UserID:              user.ID,
		Puntos:              0,
		Acciones:            0,
		TorneosParticipados: 0,
		CantidadAmigos:      0,
		EsDuenoTorneo:       false,
		TorneosGanados:      0,
		PendingMedalla:      0,
		PendingAmigo:        0,
	}

	if err := h.repo.UpdateUserStats(&stats); err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	profile := models.UserProfile{
		UserID:           user.ID,
		Slogan:           "ViVe tu mejor vida",
		Cabello:          "default",
		Vestimenta:       "default",
		Barba:            "0",
		DetalleFacial:    "0",
		DetalleAdicional: "0",
	}

	if err := h.repo.UpdateUserProfile(&profile); err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var userAccess models.CreateUserAccess
	if err := json.NewDecoder(r.Body).Decode(&userAccess); err != nil {
		http.Error(w, "Error al decodificar el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}

	existingUser, err := h.repo.GetUserByUsernameAndPassword(userAccess.Username, userAccess.Password)
	if err != nil && existingUser == nil {
		http.Error(w, "Credenciales incorrectas", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(existingUser)
}

func (h *UserHandler) CreateOrUpdateUserBasicInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	basicInfo := models.UserBasicInfo{
		UserID: id,
	}

	if err := json.NewDecoder(r.Body).Decode(&basicInfo); err != nil {
		http.Error(w, "Error al decodificar el cuerpo de la solicitud", http.StatusBadRequest)

		return
	}

	err := h.repo.CreateOrUpdateUserBasicInfo(&basicInfo)
	if err != nil {
		http.Error(w, "Error al guardar la información básica del usuario "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(basicInfo)
}

func (h *UserHandler) GetUserBasicInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	basicInfo := models.UserBasicInfo{
		UserID: id,
	}

	basicInfoUser, err := h.repo.GetUserBasicInfo(&basicInfo)
	if err != nil {
		http.Error(w, "Error, no se pudo obtener la información básica del usuario "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(basicInfoUser)
}

func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	if limit <= 0 {
		limit = 10
	}

	users, err := h.repo.ListUsers(limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	user, err := h.repo.GetUserByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) ReLoginUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	user, err := h.repo.ReLoginUserByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var user models.UserAccess
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Error al decodificar el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}

	user.ID = id
	if err := h.repo.UpdateUser(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	profile, err := h.repo.GetUserProfile(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profile)
}

func (h *UserHandler) UpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var profile models.UserProfile
	if err := json.NewDecoder(r.Body).Decode(&profile); err != nil {
		http.Error(w, "Error al decodificar el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}

	profile.UserID = id
	if err := h.repo.UpdateUserProfile(&profile); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profile)
}

func (h *UserHandler) UpdateUserProfileEdit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var profile models.EditProfile
	if err := json.NewDecoder(r.Body).Decode(&profile); err != nil {
		http.Error(w, "Error al decodificar el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}

	profile.UserID = id
	if err := h.repo.EditUserProfile(&profile); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profile)
}

func (h *UserHandler) GetUserStats(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	stats, err := h.repo.GetUserStats(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func (h *UserHandler) UpdateUserStats(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var stats models.UserStats
	if err := json.NewDecoder(r.Body).Decode(&stats); err != nil {
		http.Error(w, "Error al decodificar el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}

	stats.UserID = id
	if err := h.repo.UpdateUserStats(&stats); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func (h *UserHandler) GetRanking(w http.ResponseWriter, r *http.Request) {
	ranking, err := h.repo.GetRanking()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ranking)
}
