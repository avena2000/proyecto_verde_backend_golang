package handlers

import (
	"backend_proyecto_verde/internal/models"
	"backend_proyecto_verde/internal/repository/postgres"
	"backend_proyecto_verde/internal/utils"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
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
		utils.RespondWithBadRequest(w, "Error al decodificar el cuerpo de la solicitud", err.Error())
		return
	}

	if err := h.repo.CreateMedalla(&medalla); err != nil {
		utils.RespondWithDatabaseError(w, "Error al crear la medalla", err.Error())
		return
	}

	utils.RespondWithCreated(w, medalla, "Medalla creada correctamente")
}

func (h *MedallasHandler) GetMedallas(w http.ResponseWriter, r *http.Request) {
	medallas, err := h.repo.GetMedallas()
	if err != nil {
		utils.RespondWithDatabaseError(w, "Error al obtener las medallas", err.Error())
		return
	}

	utils.RespondWithSuccess(w, medallas, "Medallas obtenidas correctamente")
}

func (h *MedallasHandler) AsignarMedalla(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]
	medallaID := vars["medalla_id"]

	if err := h.repo.AsignarMedalla(userID, medallaID); err != nil {
		utils.RespondWithDatabaseError(w, "Error al asignar la medalla", err.Error())
		return
	}

	utils.RespondWithSuccess(w, nil, "Medalla asignada correctamente")
}

func (h *MedallasHandler) GetMedallasUsuario(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]

	medallas, err := h.repo.GetMedallasUsuario(userID)
	if err != nil {
		utils.RespondWithDatabaseError(w, "Error al obtener las medallas del usuario", err.Error())
		return
	}

	utils.RespondWithSuccess(w, medallas, "Medallas del usuario obtenidas correctamente")
}

func (h *MedallasHandler) VerifyAndUpdateMedallas(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]

	if err := h.repo.VerifyAndUpdateMedallas(userID); err != nil {
		utils.RespondWithDatabaseError(w, "Error al verificar y actualizar las medallas", err.Error())
		return
	}

	utils.RespondWithSuccess(w, nil, "Medallas verificadas y actualizadas correctamente")
}

func (h *MedallasHandler) GetSlogansMedallasGanadas(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]

	slogans, err := h.repo.GetSlogansMedallasGanadas(userID)
	if err != nil {
		utils.RespondWithDatabaseError(w, "Error al obtener los slogans de medallas ganadas", err.Error())
		return
	}

	utils.RespondWithSuccess(w, slogans, "Slogans de medallas ganadas obtenidos correctamente")
}

func (h *MedallasHandler) ResetPendingMedallas(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]

	if err := h.repo.ResetPendingMedallas(userID); err != nil {
		utils.RespondWithDatabaseError(w, "Error al reiniciar el contador de medallas pendientes", err.Error())
		return
	}

	utils.RespondWithSuccess(w, nil, "Contador de medallas pendientes reiniciado correctamente")
}
