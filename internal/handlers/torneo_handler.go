package handlers

import (
	"backend_proyecto_verde/internal/models"
	"backend_proyecto_verde/internal/repository/postgres"
	"backend_proyecto_verde/internal/utils"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
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
		utils.RespondWithBadRequest(w, "Error al decodificar el cuerpo de la solicitud", err.Error())
		return
	}

	// Verificar si el usuario puede crear un torneo
	_, message, err := h.repo.GetIfTournamentOwner(torneo.IDCreator)
	if err != nil {
		utils.RespondWithBadRequest(w, message, err.Error())
		return
	}

	if err := h.repo.CreateTorneo(&torneo); err != nil {
		utils.RespondWithDatabaseError(w, "Error al crear el torneo", err.Error())
		return
	}

	utils.RespondWithCreated(w, torneo, "Torneo creado correctamente")
}

func (h *TorneoHandler) GetTorneo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	torneo, err := h.repo.GetTorneoByID(id)
	if err != nil {
		utils.RespondWithNotFound(w, "Torneo no encontrado", "No se encontró el torneo con el ID proporcionado")
		return
	}

	utils.RespondWithSuccess(w, torneo, "Torneo obtenido correctamente")
}

func (h *TorneoHandler) GetTorneoByCodeID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	codeID := vars["code_id"]

	torneo, err := h.repo.GetTorneoByCodeID(codeID)
	if err != nil {
		utils.RespondWithNotFound(w, "Torneo no encontrado", "No se encontró el torneo con el ID proporcionado")
		return
	}

	utils.RespondWithSuccess(w, torneo, "Torneo obtenido correctamente")
}

func (h *TorneoHandler) GetTorneoAdmin(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	torneo, err := h.repo.GetTorneoByAdminID(id)
	if err != nil {
		utils.RespondWithNotFound(w, "Torneo no encontrado", "No se encontró el torneo con el ID proporcionado")
		return
	}

	utils.RespondWithSuccess(w, torneo, "Torneo obtenido correctamente")
}

func (h *TorneoHandler) UpdateTorneo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var torneo models.Torneo
	if err := json.NewDecoder(r.Body).Decode(&torneo); err != nil {
		utils.RespondWithBadRequest(w, "Error al decodificar el cuerpo de la solicitud", err.Error())
		return
	}

	torneo.ID = id
	if err := h.repo.UpdateTorneo(&torneo); err != nil {
		utils.RespondWithDatabaseError(w, "Error al actualizar el torneo", err.Error())
		return
	}

	utils.RespondWithSuccess(w, torneo, "Torneo actualizado correctamente")
}

func (h *TorneoHandler) ListTorneos(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	if limit <= 0 {
		limit = 10
	}

	torneos, err := h.repo.ListTorneos(limit, offset)
	if err != nil {
		utils.RespondWithDatabaseError(w, "Error al obtener la lista de torneos", err.Error())
		return
	}

	utils.RespondWithSuccess(w, torneos, "Lista de torneos obtenida correctamente")
}

func (h *TorneoHandler) GetTorneoStats(w http.ResponseWriter, r *http.Request) {
	// Esta función aún no está implementada en el repositorio
	// Devolvemos una respuesta vacía por ahora
	utils.RespondWithSuccess(w, []models.TorneoEstadisticas{}, "No hay estadísticas disponibles")
}

func (h *TorneoHandler) TerminarTorneo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.repo.TerminarTorneo(id); err != nil {
		utils.RespondWithDatabaseError(w, "Error al finalizar el torneo", err.Error())
		return
	}

	utils.RespondWithSuccess(w, nil, "Torneo finalizado correctamente")
}

func (h *TorneoHandler) BorrarTorneo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	message, err := h.repo.BorrarTorneoEstadisticas(id)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.RespondWithNotFound(w, "Torneo no encontrado", "No se encontró el torneo con el ID proporcionado")
		} else {
			utils.RespondWithDatabaseError(w, message, err.Error())
		}
		return
	}

	utils.RespondWithSuccess(w, nil, "Torneo eliminado correctamente")
}

// GetTorneosUsuario obtiene todos los torneos relacionados con un usuario específico,
// tanto los que ha creado como en los que ha participado
func (h *TorneoHandler) GetTorneosUsuario(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]

	torneos, err := h.repo.GetTorneosRelacionadosUsuario(userID)
	if err != nil {
		utils.RespondWithDatabaseError(w, "Error al obtener torneos del usuario", err.Error())
		return
	}

	utils.RespondWithSuccess(w, torneos, "Torneos del usuario obtenidos correctamente")
}

func (h *TorneoHandler) InscribirUsuario(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	codeID := vars["code_id"]

	var body struct {
		UserID string `json:"user_id"`
		Team   bool   `json:"team"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		utils.RespondWithBadRequest(w, "Error al decodificar el cuerpo de la solicitud", err.Error())
		return
	}

	// Verificar si el usuario puede inscribirse a un torneo
	_, message, err := h.repo.GetIfTournamentOwner(body.UserID)
	if err != nil {
		utils.RespondWithBadRequest(w, message, err.Error())
		return
	}

	if err := h.repo.InscribirUsuario(codeID, body.UserID, body.Team); err != nil {
		utils.RespondWithDatabaseError(w, "Error al inscribir el usuario", err.Error())
		return
	}

	utils.RespondWithCreated(w, nil, "Usuario inscrito correctamente")
}

// SalirTorneo maneja la solicitud para que un usuario abandone un torneo
func (h *TorneoHandler) SalirTorneo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	torneoID := vars["torneo_id"]
	userID := vars["user_id"]

	if err := h.repo.SalirTorneo(userID, torneoID); err != nil {
		utils.RespondWithDatabaseError(w, "Error al salir del torneo", err.Error())
		return
	}

	utils.RespondWithSuccess(w, nil, "Usuario ha salido del torneo correctamente")
}

// GetEquipoUsuarioTorneo maneja la solicitud para obtener el equipo de un usuario en un torneo específico
func (h *TorneoHandler) GetEquipoUsuarioTorneo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	torneoID := vars["torneo_id"]
	userID := vars["user_id"]

	equipo, err := h.repo.GetEquipoUsuarioTorneo(torneoID, userID)
	if err != nil {
		utils.RespondWithBadRequest(w, "Error al obtener el equipo del usuario", err.Error())
		return
	}

	response := struct {
		Equipo bool `json:"equipo"`
	}{
		Equipo: *equipo,
	}

	utils.RespondWithSuccess(w, response, "Equipo del usuario obtenido correctamente")
}

// FinalizeTournaments finaliza automáticamente todos los torneos cuya fecha de fin ya ha pasado
func (h *TorneoHandler) FinalizeTournaments() error {
	// Buscar todos los torneos vencidos
	creatorIDs, err := h.repo.FindExpiredTournaments()
	if err != nil {
		return fmt.Errorf("error al buscar torneos vencidos: %w", err)
	}

	// Si no hay torneos vencidos, retornar sin errores
	if len(creatorIDs) == 0 {
		return nil
	}

	// Finalizar cada torneo vencido
	var finalErrors []string
	for _, creatorID := range creatorIDs {
		if err := h.repo.TerminarTorneo(creatorID); err != nil {
			finalErrors = append(finalErrors, fmt.Sprintf("Error al finalizar torneo del creador %s: %v", creatorID, err))
		}
	}

	// Si hubo errores, reportarlos todos juntos
	if len(finalErrors) > 0 {
		return fmt.Errorf("errores al finalizar torneos: %s", strings.Join(finalErrors, "; "))
	}

	return nil
}

// UpdateTorneoFechaFin actualiza solo la fecha de fin de un torneo
func (h *TorneoHandler) UpdateTorneoFechaFin(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var body struct {
		FechaFin string `json:"fecha_fin"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		utils.RespondWithBadRequest(w, "Error al decodificar el cuerpo de la solicitud", err.Error())
		return
	}

	if body.FechaFin == "" {
		utils.RespondWithBadRequest(w, "La fecha de fin no puede estar vacía", "fecha_fin es requerida")
		return
	}

	if err := h.repo.UpdateTorneoFechaFin(id, body.FechaFin); err != nil {
		utils.RespondWithDatabaseError(w, "Error al actualizar la fecha de fin del torneo", err.Error())
		return
	}

	utils.RespondWithSuccess(w, map[string]string{"id": id, "fecha_fin": body.FechaFin}, "Fecha de fin del torneo actualizada correctamente")
}
