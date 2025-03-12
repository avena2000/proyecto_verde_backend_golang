package models

type UserAccess struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"` // El guión evita que se serialice en JSON
}

type LoginUserAccess struct {
	ID                    string `json:"id"`
	Username              string `json:"username"`
	IsPersonalInformation bool   `json:"is_personal_information"`
}

type CreateUserAccess struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UserProfile struct {
	ID               string `json:"id"`
	UserID           string `json:"user_id"`
	Slogan           string `json:"slogan"`
	Cabello          string `json:"cabello"`
	Vestimenta       string `json:"vestimenta"`
	Barba            string `json:"barba"`
	DetalleFacial    string `json:"detalle_facial"`
	DetalleAdicional string `json:"detalle_adicional"`
}

type EditProfile struct {
	UserID           string  `json:"user_id"`
	Nombre           *string `json:"nombre,omitempty"`
	Apellido         *string `json:"apellido,omitempty"`
	Slogan           *string `json:"slogan,omitempty"`
	Cabello          *string `json:"cabello,omitempty"`
	Vestimenta       *string `json:"vestimenta,omitempty"`
	Barba            *string `json:"barba,omitempty"`
	DetalleFacial    *string `json:"detalle_facial,omitempty"`
	DetalleAdicional *string `json:"detalle_adicional,omitempty"`
}

type UserBasicInfo struct {
	ID       string `json:"id"`
	UserID   string `json:"user_id"`
	Numero   string `json:"numero"`
	Nombre   string `json:"nombre"`
	Apellido string `json:"apellido"`
	FriendId string `json:"friend_id"`
}

type UserFriend struct {
	ID               string  `json:"id"`
	FriendID         string  `json:"friend_id"`
	Nombre           string  `json:"nombre"`
	Apellido         string  `json:"apellido"`
	PendingID        *string `json:"pending_id,omitempty"`
	Slogan           string  `json:"slogan"`
	Cabello          string  `json:"cabello"`
	Vestimenta       string  `json:"vestimenta"`
	Barba            string  `json:"barba"`
	DetalleFacial    string  `json:"detalle_facial"`
	DetalleAdicional string  `json:"detalle_adicional"`
}

type UserStats struct {
	ID                  string  `json:"id"`
	UserID              string  `json:"user_id"`
	Puntos              int     `json:"puntos"`
	Acciones            int     `json:"acciones"`
	TorneosParticipados int     `json:"torneos_participados"`
	TorneosGanados      int     `json:"torneos_ganados"`
	CantidadAmigos      int     `json:"cantidad_amigos"`
	EsDuenoTorneo       bool    `json:"es_dueno_torneo"`
	PendingMedalla      int     `json:"pending_medalla"`
	PendingAmigo        int     `json:"pending_amigo"`
	TorneoId            *string `json:"torneo_id"`
}
