package models

import "time"

type UserAction struct {
	ID             string     `json:"id"`
	UserID         string     `json:"user_id"`
	TipoAccion     string     `json:"tipo_accion"`
	Foto           string     `json:"foto"`
	Latitud        float64    `json:"latitud"`
	Longitud       float64    `json:"longitud"`
	EnColaboracion bool       `json:"en_colaboracion"`
	Colaboradores  []string   `json:"colaboradores,omitempty"`
	EsParaTorneo   bool       `json:"es_para_torneo"`
	IDTorneo       *string    `json:"id_torneo,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
} 