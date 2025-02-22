package models

import "time"

type Torneo struct {
	ID                  string    `json:"id"`
	IDCreator           string    `json:"id_creator"`
	Nombre              string    `json:"nombre"`
	Modalidad           string    `json:"modalidad"`
	UbicacionALatitud   float64   `json:"ubicacion_a_latitud"`
	UbicacionALongitud  float64   `json:"ubicacion_a_longitud"`
	NombreUbicacionA    string    `json:"nombre_ubicacion_a"`
	UbicacionBLatitud   *float64  `json:"ubicacion_b_latitud,omitempty"`
	UbicacionBLongitud  *float64  `json:"ubicacion_b_longitud,omitempty"`
	NombreUbicacionB    *string   `json:"nombre_ubicacion_b,omitempty"`
	FechaInicio         time.Time `json:"fecha_inicio"`
	FechaFin            time.Time `json:"fecha_fin"`
	UbicacionAproximada bool      `json:"ubicacion_aproximada"`
	MetrosAprox         *int      `json:"metros_aproximados,omitempty"`
	Finalizado          bool      `json:"finalizado"`
	CodeID              string    `json:"code_id"`
	GanadorVersus       *bool     `json:"ganador_versus,omitempty"`
	GanadorIndividual   *string   `json:"ganador_individual,omitempty"`
}

type TorneoEstadisticas struct {
	ID         string `json:"id"`
	IDJugador  string `json:"id_jugador"`
	Equipo     bool   `json:"equipo"`
	IDTorneo   string `json:"id_torneo"`
	Modalidad  string `json:"modalidad"`
	Puntos     int    `json:"puntos"`
	Habilitado bool   `json:"habilitado"`
} 