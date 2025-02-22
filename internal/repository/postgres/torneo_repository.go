package postgres

import (
	"backend_proyecto_verde/internal/models"
	"backend_proyecto_verde/internal/utils"
	"database/sql"
	"fmt"
)

type TorneoRepository struct {
	db *sql.DB
}

func NewTorneoRepository(db *sql.DB) *TorneoRepository {
	return &TorneoRepository{db: db}
}

func (r *TorneoRepository) CreateTorneo(torneo *models.Torneo) error {
	torneo.CodeID = utils.GenerateUniqueFriendId(r.db, true)

	query := `
		INSERT INTO torneos (
			nombre, id_creator, modalidad, ubicacion_a_latitud, ubicacion_a_longitud,
			nombre_ubicacion_a, ubicacion_b_latitud, ubicacion_b_longitud,
			nombre_ubicacion_b, fecha_inicio, fecha_fin, ubicacion_aproximada,
			metros_aproximados, code_id
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING id`

	err := r.db.QueryRow(
		query,
		torneo.Nombre,
		torneo.IDCreator,
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
		torneo.MetrosAprox,
		torneo.CodeID, // Se agregó el valor faltante
	).Scan(&torneo.ID)

	if err != nil {
		return fmt.Errorf("error al crear torneo: %w", err)
	}
	_, err = r.db.Exec(`UPDATE user_stats SET es_dueno_torneo = true WHERE user_id = $1`, torneo.IDCreator)
	if err != nil {
		return fmt.Errorf("error al actualizar es_dueno_torneo: %w", err)
	}
	return nil
}
func (r *TorneoRepository) GetTorneoByID(id string) (*models.Torneo, error) {
	query := `
		SELECT id, nombre, modalidad, ubicacion_a_latitud, ubicacion_a_longitud,
			nombre_ubicacion_a, ubicacion_b_latitud, ubicacion_b_longitud,
			nombre_ubicacion_b, fecha_inicio, fecha_fin, ubicacion_aproximada,
			metros_aproximados, finalizado, ganador_versus, ganador_individual
		FROM torneos
		WHERE id = $1`

	torneo := &models.Torneo{}
	err := r.db.QueryRow(query, id).Scan(
		&torneo.ID, &torneo.Nombre, &torneo.Modalidad,
		&torneo.UbicacionALatitud, &torneo.UbicacionALongitud,
		&torneo.NombreUbicacionA, &torneo.UbicacionBLatitud,
		&torneo.UbicacionBLongitud, &torneo.NombreUbicacionB,
		&torneo.FechaInicio, &torneo.FechaFin,
		&torneo.UbicacionAproximada, &torneo.MetrosAprox,
		&torneo.Finalizado, &torneo.GanadorVersus, &torneo.GanadorIndividual,
	)
	if err != nil {
		return nil, err
	}
	return torneo, nil
}

func (r *TorneoRepository) GetTorneoByAdminID(id string) (*models.Torneo, error) {
	query := `
		SELECT id, id_creator, nombre, modalidad, ubicacion_a_latitud, ubicacion_a_longitud,
			nombre_ubicacion_a, ubicacion_b_latitud, ubicacion_b_longitud,
			nombre_ubicacion_b, fecha_inicio, fecha_fin, ubicacion_aproximada,
			metros_aproximados, code_id, finalizado, ganador_versus, ganador_individual
		FROM torneos
		WHERE id_creator = $1`

	torneo := &models.Torneo{}
	err := r.db.QueryRow(query, id).Scan(
		&torneo.ID, &torneo.IDCreator, &torneo.Nombre, &torneo.Modalidad,
		&torneo.UbicacionALatitud, &torneo.UbicacionALongitud,
		&torneo.NombreUbicacionA, &torneo.UbicacionBLatitud,
		&torneo.UbicacionBLongitud, &torneo.NombreUbicacionB,
		&torneo.FechaInicio, &torneo.FechaFin,
		&torneo.UbicacionAproximada, &torneo.MetrosAprox,
		&torneo.CodeID, &torneo.Finalizado, &torneo.GanadorVersus, &torneo.GanadorIndividual,
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
			ubicacion_aproximada = $11, metros_aproximados = $12,
			finalizado = $13, ganador_versus = $14, ganador_individual = $15
		WHERE id = $16`

	_, err := r.db.Exec(query,
		torneo.Nombre, torneo.Modalidad,
		torneo.UbicacionALatitud, torneo.UbicacionALongitud,
		torneo.NombreUbicacionA, torneo.UbicacionBLatitud,
		torneo.UbicacionBLongitud, torneo.NombreUbicacionB,
		torneo.FechaInicio, torneo.FechaFin,
		torneo.UbicacionAproximada, torneo.MetrosAprox,
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

func (r *TorneoRepository) TerminarTorneo(idCreator string) error {

	var torneoID string
	err := r.db.QueryRow(`
		SELECT id FROM torneos
		WHERE id_creator = $1
	`, idCreator).Scan(&torneoID)
	if err != nil {
		return err
	}

	// Verificar que el torneo exista y pertenezca al creador
	var modalidad string
	err = r.db.QueryRow(`
		SELECT modalidad FROM torneos
		WHERE id = $1
	`, torneoID).Scan(&modalidad)
	if err != nil {
		return err
	}

	if modalidad == "Versus" {
		// Obtener las estadísticas del torneo
		rows, err := r.db.Query(`
			SELECT id_jugador, equipo, puntos
			FROM torneo_estadisticas
			WHERE id_torneo = $1 AND habilitado = true
		`, torneoID)
		if err != nil {
			return err
		}
		defer rows.Close()

		// Separar jugadores por equipo y sumar puntos
		equipoA := make(map[string]int)
		equipoB := make(map[string]int)
		puntosEquipoA := 0
		puntosEquipoB := 0

		for rows.Next() {
			var idJugador string
			var equipo bool
			var puntos int
			if err := rows.Scan(&idJugador, &equipo, &puntos); err != nil {
				return err
			}

			if equipo {
				equipoB[idJugador] = puntos
				puntosEquipoB += puntos
			} else {
				equipoA[idJugador] = puntos
				puntosEquipoA += puntos
			}
		}

		// Determinar el equipo ganador
		equipoGanador := puntosEquipoA > puntosEquipoB

		// Actualizar el torneo con el ganador
		_, err = r.db.Exec(`
			UPDATE torneos
			SET finalizado = true, ganador_versus = $1
			WHERE id = $2
		`, equipoGanador, torneoID)
		if err != nil {
			return err
		}

		// Actualizar estadísticas de los jugadores
		tx, err := r.db.Begin()
		if err != nil {
			return err
		}
		defer tx.Rollback()

		// Actualizar torneos_participados para todos los jugadores
		for idJugador := range equipoA {
			_, err = tx.Exec(`
				UPDATE user_stats
				SET torneos_participados = torneos_participados + 1
				WHERE user_id = $1
			`, idJugador)
			if err != nil {
				return err
			}
		}
		for idJugador := range equipoB {
			_, err = tx.Exec(`
				UPDATE user_stats
				SET torneos_participados = torneos_participados + 1
				WHERE user_id = $1
			`, idJugador)
			if err != nil {
				return err
			}
		}

		// Actualizar torneos_ganados para el equipo ganador
		equipoGanadorMap := equipoA
		if !equipoGanador {
			equipoGanadorMap = equipoB
		}
		for idJugador := range equipoGanadorMap {
			_, err = tx.Exec(`
				UPDATE user_stats
				SET torneos_ganados = torneos_ganados + 1
				WHERE user_id = $1
			`, idJugador)
			if err != nil {
				return err
			}
		}

		// Deshabilitar todas las estadísticas del torneo
		_, err = tx.Exec(`
			UPDATE torneo_estadisticas
			SET habilitado = false
			WHERE id_torneo = $1
		`, torneoID)
		if err != nil {
			return err
		}

		return tx.Commit()
	}

	return nil
}

func (r *TorneoRepository) BorrarTorneoEstadisticas(idCreator string) error {

	var torneoID string
	err := r.db.QueryRow(`
		SELECT id FROM torneos
		WHERE id_creator = $1
	`, idCreator).Scan(&torneoID)
	if err != nil {
		return err
	}

	// Verificar que el torneo exista y pertenezca al creador
	var exists bool
	err = r.db.QueryRow(`
		SELECT EXISTS(SELECT 1 FROM torneos WHERE id = $1 AND id_creator = $2)
	`, torneoID, idCreator).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return sql.ErrNoRows
	}

	_, err = r.db.Exec(`
		UPDATE torneos
		SET finalizado = true
		WHERE id = $1
	`, torneoID)
	if err != nil {
		return err
	}

	// Borrar todas las estadísticas del torneo
	_, err = r.db.Exec(`
		DELETE FROM torneo_estadisticas
		WHERE id_torneo = $1
	`, torneoID)

	if err != nil {
		return err
	}
	_, err = r.db.Exec(`
		UPDATE user_stats
		SET es_dueno_torneo = false
		WHERE user_id = $1
	`, idCreator)

	return err
}

func (r *TorneoRepository) InscribirUsuario(codeID string, userID string) error {
	// Obtener el torneo por code_id y verificar que esté habilitado
	var torneoID string
	var finalizado bool
	var modalidad string
	err := r.db.QueryRow(`
		SELECT id, finalizado, modalidad 
		FROM torneos 
		WHERE code_id = $1
	`, codeID).Scan(&torneoID, &finalizado, &modalidad)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("torneo no encontrado")
		}
		return err
	}

	if finalizado {
		return fmt.Errorf("el torneo ya está finalizado")
	}

	// Verificar si el usuario ya está inscrito
	var exists bool
	err = r.db.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM torneo_estadisticas 
			WHERE id_torneo = $1 AND id_jugador = $2
		)
	`, torneoID, userID).Scan(&exists)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("el usuario ya está inscrito en este torneo")
	}

	// Iniciar transacción
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insertar en torneo_estadisticas
	_, err = tx.Exec(`
		INSERT INTO torneo_estadisticas (
			id_jugador, id_torneo, modalidad, puntos, 
			equipo, habilitado
		) VALUES ($1, $2, $3, 0, false, true)
	`, userID, torneoID, modalidad)
	if err != nil {
		return err
	}

	// Actualizar user_stats con el id del torneo
	_, err = tx.Exec(`
		UPDATE user_stats 
		SET torneo_id = $1 
		WHERE user_id = $2
	`, torneoID, userID)
	if err != nil {
		return err
	}

	return tx.Commit()
}
