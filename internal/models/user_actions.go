package models

import "time"

type UserAction struct {
	ID             string     `json:"id"`
	UserID         string     `json:"user_id"`
	TipoAccion     string     `json:"tipo_accion"`
	Foto           string     `json:"foto"`
	Latitud        float64    `json:"latitud"`
	Longitud       float64    `json:"longitud"`
	Ciudad         string     `json:"ciudad"`
	Lugar          string     `json:"lugar"`
	EnColaboracion bool       `json:"en_colaboracion"`
	Colaboradores  *[]string  `json:"colaboradores,omitempty"`
	EsParaTorneo   bool       `json:"es_para_torneo"`
	IDTorneo       *string    `json:"id_torneo,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
}

type UserRanking struct {
	UserID           string `json:"user_id"`
	Puntos           int    `json:"puntos"`
	Acciones         int    `json:"acciones"`
	TorneosGanados   int    `json:"torneos_ganados"`
	CantidadAmigos   int    `json:"cantidad_amigos"`
	Slogan           string `json:"slogan"`
	Cabello          string `json:"cabello"`
	Vestimenta       string `json:"vestimenta"`
	Barba            string `json:"barba"`
	DetalleFacial    string `json:"detalle_facial"`
	DetalleAdicional string `json:"detalle_adicional"`
	Nombre           string `json:"nombre"`
	Apellido         string `json:"apellido"`
}
