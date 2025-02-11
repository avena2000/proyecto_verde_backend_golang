package models

import "time"

type Medalla struct {
	ID                    string `json:"id"`
	Nombre                string `json:"nombre"`
	Descripcion           string `json:"descripcion"`
	Dificultad            int    `json:"dificultad"`
	RequiereAmistades     bool   `json:"requiere_amistades"`
	RequierePuntos        bool   `json:"requiere_puntos"`
	RequiereAcciones      bool   `json:"requiere_acciones"`
	RequiereTorneos       bool   `json:"requiere_torneos"`
	RequiereVictoriaTorneos bool `json:"requiere_victoria_torneos"`
	NumeroRequerido       *int   `json:"numero_requerido,omitempty"`
}

type MedallaGanada struct {
	ID          string    `json:"id"`
	IDUsuario   string    `json:"id_usuario"`
	IDMedalla   string    `json:"id_medalla"`
	FechaGanada time.Time `json:"fecha_ganada"`
} 