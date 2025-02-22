package handlers

import (
	"backend_proyecto_verde/internal/repository/postgres"
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
		http.Error(w, "Error al decodificar el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}

	if err := h.repo.SendFriendRequest(userID, friendIDRequest.ID); err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "El usuario no existe", http.StatusNotFound)
			return
		}
		if errors.Is(err, postgres.ErrSelfFriendRequest) {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		if errors.Is(err, postgres.ErrFriendRequestExists) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *UserFriendsHandler) AcceptFriendRequest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]
	friendID := vars["friend_id"]

	if err := h.repo.AcceptFriendRequest(userID, friendID); err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Solicitud de amistad no encontrada", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *UserFriendsHandler) GetFriendsList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]

	friends, err := h.repo.GetFriendsList(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(friends)
}

func (h *UserFriendsHandler) RemoveFriend(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]
	friendID := vars["friend_id"]

	if err := h.repo.RemoveFriend(userID, friendID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
