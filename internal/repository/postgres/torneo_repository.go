package postgres

import (
	"backend_proyecto_verde/internal/models"
	"database/sql"
)

type TorneoRepository struct {
	db *sql.DB
}

func NewTorneoRepository(db *sql.DB) *TorneoRepository {
	return &TorneoRepository{db: db}
}

func (r *TorneoRepository) CreateTorneo(torneo *models.Torneo) error {
	query := `
		INSERT INTO torneos (
			nombre, modalidad, ubicacion_a_latitud, ubicacion_a_longitud,
			nombre_ubicacion_a, ubicacion_b_latitud, ubicacion_b_longitud,
			nombre_ubicacion_b, fecha_inicio, fecha_fin, ubicacion_aproximada,
			kilometros_aproximados
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id`

	return r.db.QueryRow(
		query,
		torneo.Nombre,
		torneo.Modalidad,
		torneo.UbicacionALatitud,
		torneo.UbicacionALongitud,
		torneo.NombreUbicacionA,
		torneo.UbicacionBLatitud,
		torneo.UbicacionBLongitud,
		torneo.NombreUbicacionB,
		torneo.FechaInicio,
		torneo.FechaFin,
		torneo.UbicacionAproximada,
		torneo.KilometrosAprox,
	).Scan(&torneo.ID)
}

func (r *TorneoRepository) GetTorneoByID(id string) (*models.Torneo, error) {
	query := `
		SELECT id, nombre, modalidad, ubicacion_a_latitud, ubicacion_a_longitud,
			nombre_ubicacion_a, ubicacion_b_latitud, ubicacion_b_longitud,
			nombre_ubicacion_b, fecha_inicio, fecha_fin, ubicacion_aproximada,
			kilometros_aproximados, finalizado, ganador_versus, ganador_individual
		FROM torneos
		WHERE id = $1`

	torneo := &models.Torneo{}
	err := r.db.QueryRow(query, id).Scan(
		&torneo.ID, &torneo.Nombre, &torneo.Modalidad,
		&torneo.UbicacionALatitud, &torneo.UbicacionALongitud,
		&torneo.NombreUbicacionA, &torneo.UbicacionBLatitud,
		&torneo.UbicacionBLongitud, &torneo.NombreUbicacionB,
		&torneo.FechaInicio, &torneo.FechaFin,
		&torneo.UbicacionAproximada, &torneo.KilometrosAprox,
		&torneo.Finalizado, &torneo.GanadorVersus, &torneo.GanadorIndividual,
	)
	if err != nil {
		return nil, err
	}
	return torneo, nil
}

func (r *TorneoRepository) UpdateTorneo(torneo *models.Torneo) error {
	query := `
		UPDATE torneos
		SET nombre = $1, modalidad = $2, ubicacion_a_latitud = $3,
			ubicacion_a_longitud = $4, nombre_ubicacion_a = $5,
			ubicacion_b_latitud = $6, ubicacion_b_longitud = $7,
			nombre_ubicacion_b = $8, fecha_inicio = $9, fecha_fin = $10,
			ubicacion_aproximada = $11, kilometros_aproximados = $12,
			finalizado = $13, ganador_versus = $14, ganador_individual = $15
		WHERE id = $16`

	_, err := r.db.Exec(query,
		torneo.Nombre, torneo.Modalidad,
		torneo.UbicacionALatitud, torneo.UbicacionALongitud,
		torneo.NombreUbicacionA, torneo.UbicacionBLatitud,
		torneo.UbicacionBLongitud, torneo.NombreUbicacionB,
		torneo.FechaInicio, torneo.FechaFin,
		torneo.UbicacionAproximada, torneo.KilometrosAprox,
		torneo.Finalizado, torneo.GanadorVersus,
		torneo.GanadorIndividual, torneo.ID)
	return err
}

func (r *TorneoRepository) ListTorneos(limit, offset int) ([]models.Torneo, error) {
	query := `
		SELECT id, nombre, modalidad, ubicacion_a_latitud, ubicacion_a_longitud,
			nombre_ubicacion_a, fecha_inicio, fecha_fin, finalizado
		FROM torneos
		ORDER BY fecha_inicio DESC
		LIMIT $1 OFFSET $2`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var torneos []models.Torneo
	for rows.Next() {
		var t models.Torneo
		err := rows.Scan(
			&t.ID, &t.Nombre, &t.Modalidad,
			&t.UbicacionALatitud, &t.UbicacionALongitud,
			&t.NombreUbicacionA, &t.FechaInicio,
			&t.FechaFin, &t.Finalizado,
		)
		if err != nil {
			return nil, err
		}
		torneos = append(torneos, t)
	}
	return torneos, nil
}
