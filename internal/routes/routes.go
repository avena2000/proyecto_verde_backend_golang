package routes

import (
	"backend_proyecto_verde/internal/handlers"
	"backend_proyecto_verde/internal/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

func SetupRoutes(
	userHandler *handlers.UserHandler,
	torneoHandler *handlers.TorneoHandler,
	userActionsHandler *handlers.UserActionsHandler,
	userFriendsHandler *handlers.UserFriendsHandler,
	medallasHandler *handlers.MedallasHandler,
) *mux.Router {
	r := mux.NewRouter()

	// Agregar middleware de logging
	r.Use(middleware.LoggingMiddleware)

	// Health check endpoint
	r.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	// Endpoint para verificar que las rutas /v1 funcionan
	r.HandleFunc("/v1/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("API v1 está funcionando correctamente"))
	}).Methods("GET")

	// Rutas de usuario
	r.HandleFunc("/api/users", userHandler.CreateUser).Methods("POST")
	r.HandleFunc("/api/auth/login", userHandler.LoginUser).Methods("POST")
	r.HandleFunc("/api/auth/relogin/{id}", userHandler.ReLoginUser).Methods("GET")
	r.HandleFunc("/api/users/{id}", userHandler.GetUser).Methods("GET")
	r.HandleFunc("/api/users/{id}", userHandler.UpdateUser).Methods("PUT")

	r.HandleFunc("/api/users/{id}/basic-info", userHandler.CreateOrUpdateUserBasicInfo).Methods("POST")
	r.HandleFunc("/api/users/{id}/basic-info", userHandler.GetUserBasicInfo).Methods("GET")
	r.HandleFunc("/api/users/{id}/basic-info", userHandler.CreateOrUpdateUserBasicInfo).Methods("PUT")
	r.HandleFunc("/api/users/{id}/profile", userHandler.GetUserProfile).Methods("GET")
	r.HandleFunc("/api/users/{id}/profile/edit", userHandler.UpdateUserProfileEdit).Methods("PUT")
	r.HandleFunc("/api/users/{id}/stats", userHandler.GetUserStats).Methods("GET")
	r.HandleFunc("/api/users/{id}/stats", userHandler.UpdateUserStats).Methods("PUT")

	// Rutas de ranking
	r.HandleFunc("/api/ranking", userHandler.GetRanking).Methods("GET")
	r.HandleFunc("/api/ranking/torneo/{torneo_id}", userHandler.GetRankingTorneo).Methods("GET")

	// Rutas de torneos
	r.HandleFunc("/api/torneos", torneoHandler.CreateTorneo).Methods("POST")
	r.HandleFunc("/api/torneos", torneoHandler.ListTorneos).Methods("GET")
	r.HandleFunc("/api/torneos/{id}", torneoHandler.GetTorneo).Methods("GET")
	r.HandleFunc("/api/torneos/code/{code_id}", torneoHandler.GetTorneoByCodeID).Methods("GET")
	r.HandleFunc("/api/torneos/admin/{id}", torneoHandler.GetTorneoAdmin).Methods("GET")
	r.HandleFunc("/api/torneos/admin/{id}/terminar", torneoHandler.TerminarTorneo).Methods("POST")
	r.HandleFunc("/api/torneos/admin/{id}/borrar", torneoHandler.BorrarTorneo).Methods("POST")
	r.HandleFunc("/api/torneos/inscribir/{code_id}", torneoHandler.InscribirUsuario).Methods("POST")
	r.HandleFunc("/api/torneos/{torneo_id}/usuario/{user_id}", torneoHandler.SalirTorneo).Methods("DELETE")
	r.HandleFunc("/api/torneos/{id}", torneoHandler.UpdateTorneo).Methods("PUT")
	r.HandleFunc("/api/torneos/{id}/fecha_fin", torneoHandler.UpdateTorneoFechaFin).Methods("PUT")
	r.HandleFunc("/api/torneos/{id}/estadisticas", torneoHandler.GetTorneoStats).Methods("GET")
	r.HandleFunc("/api/users/{user_id}/torneos", torneoHandler.GetTorneosUsuario).Methods("GET")
	r.HandleFunc("/api/torneos/{torneo_id}/usuario/{user_id}/equipo", torneoHandler.GetEquipoUsuarioTorneo).Methods("GET")

	// Rutas de acciones de usuario
	r.HandleFunc("/api/users/{user_id}/actions", userActionsHandler.CreateAction).Methods("POST")
	r.HandleFunc("/api/users/{user_id}/actions", userActionsHandler.GetUserActions).Methods("GET")
	r.HandleFunc("/api/actions/{id}", userActionsHandler.DeleteAction).Methods("DELETE")
	r.HandleFunc("/api/actions", userActionsHandler.GetAllActions).Methods("GET")

	// Rutas de amigos
	r.HandleFunc("/api/users/{user_id}/friends", userFriendsHandler.GetFriendsList).Methods("GET")
	r.HandleFunc("/api/users/{user_id}/friends/add", userFriendsHandler.SendFriendRequest).Methods("POST")
	r.HandleFunc("/api/users/{user_id}/friends/{friend_id}/accept", userFriendsHandler.AcceptFriendRequest).Methods("PUT")
	r.HandleFunc("/api/users/{user_id}/friends/{friend_id}", userFriendsHandler.RemoveFriend).Methods("DELETE")

	// Rutas de medallas
	r.HandleFunc("/api/medallas", medallasHandler.CreateMedalla).Methods("POST")
	r.HandleFunc("/api/medallas", medallasHandler.GetMedallas).Methods("GET")
	r.HandleFunc("/api/users/{user_id}/medallas", medallasHandler.GetMedallasUsuario).Methods("GET")
	r.HandleFunc("/api/users/{user_id}/medallas/{medalla_id}", medallasHandler.AsignarMedalla).Methods("POST")
	r.HandleFunc("/api/users/{user_id}/medallas/slogans", medallasHandler.GetSlogansMedallasGanadas).Methods("GET")
	r.HandleFunc("/api/users/{user_id}/medallas/reset-pending", medallasHandler.ResetPendingMedallas).Methods("GET")

	//Ruta proxy

	return r
}
