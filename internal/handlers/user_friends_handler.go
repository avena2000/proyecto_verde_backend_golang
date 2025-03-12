package handlers

import (
	"backend_proyecto_verde/internal/repository/postgres"
	"backend_proyecto_verde/internal/utils"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
)

type UserFriendsHandler struct {
	repo *postgres.UserFriendsRepository
}

func NewUserFriendsHandler(repo *postgres.UserFriendsRepository) *UserFriendsHandler {
	return &UserFriendsHandler{repo: repo}
}

func (h *UserFriendsHandler) SendFriendRequest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]

	var friendIDRequest struct {
		ID string `json:"friend_id_request"`
	}

	if err := json.NewDecoder(r.Body).Decode(&friendIDRequest); err != nil {
		utils.RespondWithBadRequest(w, "Error al decodificar el cuerpo de la solicitud", err.Error())
		return
	}

	if err := h.repo.SendFriendRequest(userID, friendIDRequest.ID); err != nil {
		if err == sql.ErrNoRows {
			utils.RespondWithNotFound(w, "El usuario no existe", "No se encontró el usuario con el ID proporcionado")
			return
		}
		if errors.Is(err, postgres.ErrSelfFriendRequest) {
			utils.RespondWithConflict(w, "No se puede enviar una solicitud de amistad a uno mismo", "No se puede enviar una solicitud de amistad a uno mismo")
			return
		}
		if errors.Is(err, postgres.ErrFriendRequestExists) {
			utils.RespondWithBadRequest(w, "Ya existe una solicitud de amistad", "Ya existe una solicitud de amistad con el ID proporcionado")
			return
		}
		utils.RespondWithDatabaseError(w, "Error al enviar la solicitud de amistad", err.Error())
		return
	}

	utils.RespondWithCreated(w, nil, "Solicitud de amistad enviada correctamente")
}

func (h *UserFriendsHandler) AcceptFriendRequest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]
	friendID := vars["friend_id"]

	if err := h.repo.AcceptFriendRequest(userID, friendID); err != nil {
		if err == sql.ErrNoRows {
			utils.RespondWithNotFound(w, "Solicitud de amistad no encontrada", "No se encontró la solicitud de amistad con el ID proporcionado")
			return
		}
		utils.RespondWithDatabaseError(w, "Error al aceptar la solicitud de amistad", err.Error())
		return
	}

	utils.RespondWithSuccess(w, nil, "Solicitud de amistad aceptada correctamente")
}

func (h *UserFriendsHandler) GetFriendsList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]

	friends, err := h.repo.GetFriendsList(userID)
	if err != nil {
		utils.RespondWithDatabaseError(w, "Error al obtener la lista de amigos", err.Error())
		return
	}

	utils.RespondWithSuccess(w, friends, "Lista de amigos obtenida correctamente")
}

func (h *UserFriendsHandler) RemoveFriend(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]
	friendID := vars["friend_id"]

	if err := h.repo.RemoveFriend(userID, friendID); err != nil {
		utils.RespondWithDatabaseError(w, "Error al eliminar la amistad", err.Error())
		return
	}

	utils.RespondWithSuccess(w, nil, "Amistad eliminada correctamente")
}
