package handlers

import (
	"backend_proyecto_verde/internal/models"
	"backend_proyecto_verde/internal/repository/postgres"
	"backend_proyecto_verde/internal/utils"
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
		utils.RespondWithBadRequest(w, "Error al decodificar el cuerpo de la solicitud", err.Error())
		return
	}

	existingUser, err := h.repo.GetUserByUsername(createUser.Username)
	if err == nil && existingUser != nil {
		utils.RespondWithUserAlreadyExists(w, "El usuario ya existe", "El usuario ya existe en la base de datos")
		return
	}

	user, err := h.repo.CreateUser(&createUser)
	if err != nil {
		utils.RespondWithDatabaseError(w, "Error al crear el usuario", err.Error())
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
		utils.RespondWithDatabaseError(w, "Error al actualizar las estadísticas del usuario", err.Error())
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

	if err := h.repo.CreateUserProfile(&profile); err != nil {
		utils.RespondWithDatabaseError(w, "Error al actualizar el perfil del usuario", err.Error())
		return
	}

	utils.RespondWithCreated(w, user, "Usuario creado correctamente")
}

func (h *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var userAccess models.CreateUserAccess
	if err := json.NewDecoder(r.Body).Decode(&userAccess); err != nil {
		utils.RespondWithBadRequest(w, "Error al decodificar el cuerpo de la solicitud", err.Error())
		return
	}

	existingUser, err := h.repo.GetUserByUsernameAndPassword(userAccess.Username, userAccess.Password)
	if err != nil && existingUser == nil {
		utils.RespondWithInvalidCredentials(w, "Credenciales incorrectas", "Las credenciales proporcionadas son incorrectas")
		return
	}

	utils.RespondWithSuccess(w, existingUser, "Inicio de sesión exitoso")
}

func (h *UserHandler) CreateOrUpdateUserBasicInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	basicInfo := models.UserBasicInfo{
		UserID: id,
	}

	if err := json.NewDecoder(r.Body).Decode(&basicInfo); err != nil {
		utils.RespondWithBadRequest(w, "Error al decodificar el cuerpo de la solicitud", err.Error())
		return
	}

	err := h.repo.CreateOrUpdateUserBasicInfo(&basicInfo)
	if err != nil {
		utils.RespondWithDatabaseError(w, "Error al guardar la información básica del usuario", err.Error())
		return
	}

	utils.RespondWithCreated(w, basicInfo, "Información básica del usuario guardada correctamente")
}

func (h *UserHandler) GetUserBasicInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	basicInfo := models.UserBasicInfo{
		UserID: id,
	}

	basicInfoUser, err := h.repo.GetUserBasicInfo(&basicInfo)
	if err != nil {
		utils.RespondWithDatabaseError(w, "Error, no se pudo obtener la información básica del usuario", err.Error())
		return
	}

	utils.RespondWithSuccess(w, basicInfoUser, "Información básica del usuario obtenida correctamente")
}

func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	if limit <= 0 {
		limit = 10
	}

	users, err := h.repo.ListUsers(limit, offset)
	if err != nil {
		utils.RespondWithDatabaseError(w, "Error al obtener la lista de usuarios", err.Error())
		return
	}

	utils.RespondWithSuccess(w, users, "Lista de usuarios obtenida correctamente")
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	user, err := h.repo.GetUserByID(id)
	if err != nil {
		utils.RespondWithNotFound(w, "Usuario no encontrado", "No se encontró el usuario con el ID proporcionado")
		return
	}

	utils.RespondWithSuccess(w, user, "Usuario obtenido correctamente")
}

func (h *UserHandler) ReLoginUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	user, err := h.repo.ReLoginUserByID(id)
	if err != nil {
		utils.RespondWithNotFound(w, "Usuario no encontrado", "No se encontró el usuario con el ID proporcionado")
		return
	}

	utils.RespondWithSuccess(w, user, "Usuario reautenticado correctamente")
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var user models.UserAccess
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		utils.RespondWithBadRequest(w, "Error al decodificar el cuerpo de la solicitud", err.Error())
		return
	}

	user.ID = id
	if err := h.repo.UpdateUser(&user); err != nil {
		utils.RespondWithDatabaseError(w, "Error al actualizar el usuario", err.Error())
		return
	}

	utils.RespondWithSuccess(w, user, "Usuario actualizado correctamente")
}

func (h *UserHandler) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	profile, err := h.repo.GetUserProfile(id)
	if err != nil {
		utils.RespondWithNotFound(w, "Perfil de usuario no encontrado", "No se encontró el perfil de usuario con el ID proporcionado")
		return
	}

	utils.RespondWithSuccess(w, profile, "Perfil de usuario obtenido correctamente")
}

func (h *UserHandler) UpdateUserProfileEdit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var profile models.EditProfile
	if err := json.NewDecoder(r.Body).Decode(&profile); err != nil {
		utils.RespondWithBadRequest(w, "Error al decodificar el cuerpo de la solicitud", err.Error())
		return
	}

	profile.UserID = id
	if err := h.repo.EditUserProfile(&profile); err != nil {
		utils.RespondWithDatabaseError(w, "Error al actualizar el perfil del usuario", err.Error())
		return
	}

	utils.RespondWithSuccess(w, profile, "Perfil de usuario actualizado correctamente")
}

func (h *UserHandler) GetUserStats(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	stats, err := h.repo.GetUserStats(id)
	if err != nil {
		utils.RespondWithNotFound(w, "Estadísticas de usuario no encontradas", "No se encontraron las estadísticas de usuario con el ID proporcionado")
		return
	}

	utils.RespondWithSuccess(w, stats, "Estadísticas de usuario obtenidas correctamente")
}

func (h *UserHandler) UpdateUserStats(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var stats models.UserStats
	if err := json.NewDecoder(r.Body).Decode(&stats); err != nil {
		utils.RespondWithBadRequest(w, "Error al decodificar el cuerpo de la solicitud", err.Error())
		return
	}

	stats.UserID = id
	if err := h.repo.UpdateUserStats(&stats); err != nil {
		utils.RespondWithDatabaseError(w, "Error al actualizar las estadísticas del usuario", err.Error())
		return
	}

	utils.RespondWithSuccess(w, stats, "Estadísticas de usuario actualizadas correctamente")
}

func (h *UserHandler) GetRanking(w http.ResponseWriter, r *http.Request) {
	ranking, err := h.repo.GetRanking()
	if err != nil {
		utils.RespondWithDatabaseError(w, "Error al obtener el ranking", err.Error())
		return
	}

	utils.RespondWithSuccess(w, ranking, "Ranking obtenido correctamente")
}
