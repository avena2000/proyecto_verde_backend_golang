package postgres

import (
	"backend_proyecto_verde/internal/models"
	"backend_proyecto_verde/internal/utils"
	"backend_proyecto_verde/pkg/database"
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

	return database.WithTransaction(r.db, func(tx *sql.Tx) error {
		query := `
			INSERT INTO torneos (
				nombre, id_creator, modalidad, ubicacion_a_latitud, ubicacion_a_longitud,
				nombre_ubicacion_a, ubicacion_b_latitud, ubicacion_b_longitud,
				nombre_ubicacion_b, fecha_inicio, fecha_fin, ubicacion_aproximada,
				metros_aproximados, code_id
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
			RETURNING id`

		err := tx.QueryRow(
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
			torneo.CodeID,
		).Scan(&torneo.ID)

		if err != nil {
			return fmt.Errorf("error al crear torneo: %w", err)
		}

		return nil
	})
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
	return database.WithTransaction(r.db, func(tx *sql.Tx) error {
		query := `
			UPDATE torneos
			SET nombre = $1, modalidad = $2, ubicacion_a_latitud = $3, ubicacion_a_longitud = $4,
				nombre_ubicacion_a = $5, ubicacion_b_latitud = $6, ubicacion_b_longitud = $7,
				nombre_ubicacion_b = $8, fecha_inicio = $9, fecha_fin = $10, ubicacion_aproximada = $11,
				metros_aproximados = $12
			WHERE id = $13`

		_, err := tx.Exec(
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
			torneo.MetrosAprox,
			torneo.ID,
		)

		if err != nil {
			return fmt.Errorf("error al actualizar torneo: %w", err)
		}

		return nil
	})
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
	return database.WithTransaction(r.db, func(tx *sql.Tx) error {
		// Obtener el torneo
		query := `
			SELECT id, modalidad
			FROM torneos
			WHERE id_creator = $1 AND finalizado = false`

		var torneoID, modalidad string
		err := tx.QueryRow(query, idCreator).Scan(&torneoID, &modalidad)
		if err != nil {
			if err == sql.ErrNoRows {
				return fmt.Errorf("no se encontró un torneo activo para este usuario")
			}
			return fmt.Errorf("error al buscar torneo: %w", err)
		}

		// Marcar el torneo como finalizado
		updateQuery := `
			UPDATE torneos
			SET finalizado = true
			WHERE id = $1`

		_, err = tx.Exec(updateQuery, torneoID)
		if err != nil {
			return fmt.Errorf("error al finalizar torneo: %w", err)
		}

		// Si es modalidad "Versus", determinar el ganador
		if modalidad == "Versus" {
			// Obtener puntos de cada equipo
			statsQuery := `
				SELECT equipo, SUM(puntos) as total_puntos
				FROM torneo_estadisticas
				WHERE id_torneo = $1
				GROUP BY equipo`

			rows, err := tx.Query(statsQuery, torneoID)
			if err != nil {
				return fmt.Errorf("error al obtener estadísticas: %w", err)
			}
			defer rows.Close()

			var equipoA, equipoB bool
			var puntosA, puntosB int

			for rows.Next() {
				var equipo bool
				var puntos int
				if err := rows.Scan(&equipo, &puntos); err != nil {
					return fmt.Errorf("error al leer estadísticas: %w", err)
				}

				if equipo {
					equipoA = true
					puntosA = puntos
				} else {
					equipoB = true
					puntosB = puntos
				}
			}

			// Determinar ganador si ambos equipos tienen participantes
			if equipoA && equipoB {
				var ganador bool
				if puntosA > puntosB {
					ganador = true
				} else {
					ganador = false
				}

				// Actualizar ganador
				updateGanadorQuery := `
					UPDATE torneos
					SET ganador_versus = $1
					WHERE id = $2`

				_, err = tx.Exec(updateGanadorQuery, ganador, torneoID)
				if err != nil {
					return fmt.Errorf("error al actualizar ganador: %w", err)
				}

				// Actualizar estadísticas de usuarios ganadores
				updateStatsQuery := `
					UPDATE user_stats
					SET torneos_ganados = torneos_ganados + 1
					FROM torneo_estadisticas
					WHERE torneo_estadisticas.id_jugador = user_stats.user_id
					AND torneo_estadisticas.id_torneo = $1
					AND torneo_estadisticas.equipo = $2`

				_, err = tx.Exec(updateStatsQuery, torneoID, ganador)
				if err != nil {
					return fmt.Errorf("error al actualizar estadísticas de usuarios: %w", err)
				}
			}
		} else if modalidad == "Individual" {
			// Para modalidad individual, determinar el ganador por puntos
			statsQuery := `
				SELECT id_jugador, puntos
				FROM torneo_estadisticas
				WHERE id_torneo = $1
				ORDER BY puntos DESC
				LIMIT 1`

			var ganadorID string
			var puntos int
			err = tx.QueryRow(statsQuery, torneoID).Scan(&ganadorID, &puntos)
			if err != nil && err != sql.ErrNoRows {
				return fmt.Errorf("error al obtener ganador: %w", err)
			}

			if err != sql.ErrNoRows {
				// Actualizar ganador
				updateGanadorQuery := `
					UPDATE torneos
					SET ganador_individual = $1
					WHERE id = $2`

				_, err = tx.Exec(updateGanadorQuery, ganadorID, torneoID)
				if err != nil {
					return fmt.Errorf("error al actualizar ganador: %w", err)
				}

				// Actualizar estadísticas del usuario ganador
				updateStatsQuery := `
					UPDATE user_stats
					SET torneos_ganados = torneos_ganados + 1
					WHERE user_id = $1`

				_, err = tx.Exec(updateStatsQuery, ganadorID)
				if err != nil {
					return fmt.Errorf("error al actualizar estadísticas del ganador: %w", err)
				}
			}
		}

		return nil
	})
}

func (r *TorneoRepository) BorrarTorneoEstadisticas(idCreator string) error {
	return database.WithTransaction(r.db, func(tx *sql.Tx) error {
		// Obtener el ID del torneo
		query := `
			SELECT id
			FROM torneos
			WHERE id_creator = $1 AND finalizado = false`

		var torneoID string
		err := tx.QueryRow(query, idCreator).Scan(&torneoID)
		if err != nil {
			if err == sql.ErrNoRows {
				return fmt.Errorf("no se encontró un torneo activo para este usuario")
			}
			return fmt.Errorf("error al buscar torneo: %w", err)
		}

		// Eliminar estadísticas del torneo
		deleteQuery := `
			DELETE FROM torneo_estadisticas
			WHERE id_torneo = $1`

		_, err = tx.Exec(deleteQuery, torneoID)
		if err != nil {
			return fmt.Errorf("error al eliminar estadísticas: %w", err)
		}

		// Eliminar el torneo
		deleteTorneoQuery := `
			DELETE FROM torneos
			WHERE id = $1`

		_, err = tx.Exec(deleteTorneoQuery, torneoID)
		if err != nil {
			return fmt.Errorf("error al eliminar torneo: %w", err)
		}

		return nil
	})
}

func (r *TorneoRepository) InscribirUsuario(codeID string, userID string) error {
	return database.WithTransaction(r.db, func(tx *sql.Tx) error {
		// Verificar si el torneo existe y no está finalizado
		torneoQuery := `
			SELECT id, modalidad
			FROM torneos
			WHERE code_id = $1 AND finalizado = false`

		var torneoID, modalidad string
		err := tx.QueryRow(torneoQuery, codeID).Scan(&torneoID, &modalidad)
		if err != nil {
			if err == sql.ErrNoRows {
				return fmt.Errorf("torneo no encontrado o ya finalizado")
			}
			return fmt.Errorf("error al buscar torneo: %w", err)
		}

		// Verificar si el usuario ya está inscrito
		checkQuery := `
			SELECT EXISTS (
				SELECT 1 FROM torneo_estadisticas
				WHERE id_torneo = $1 AND id_jugador = $2
			)`

		var exists bool
		err = tx.QueryRow(checkQuery, torneoID, userID).Scan(&exists)
		if err != nil {
			return fmt.Errorf("error al verificar inscripción: %w", err)
		}

		if exists {
			return fmt.Errorf("el usuario ya está inscrito en este torneo")
		}

		// Determinar el equipo para modalidad "Versus"
		var equipo bool
		if modalidad == "Versus" {
			// Contar usuarios en cada equipo
			countQuery := `
				SELECT equipo, COUNT(*) as total
				FROM torneo_estadisticas
				WHERE id_torneo = $1
				GROUP BY equipo`

			rows, err := tx.Query(countQuery, torneoID)
			if err != nil {
				return fmt.Errorf("error al contar equipos: %w", err)
			}
			defer rows.Close()

			var equipoACount, equipoBCount int
			for rows.Next() {
				var eq bool
				var count int
				if err := rows.Scan(&eq, &count); err != nil {
					return fmt.Errorf("error al leer conteo: %w", err)
				}

				if eq {
					equipoACount = count
				} else {
					equipoBCount = count
				}
			}

			// Asignar al equipo con menos jugadores
			equipo = equipoACount <= equipoBCount
		}

		// Inscribir al usuario
		insertQuery := `
			INSERT INTO torneo_estadisticas (id_jugador, equipo, id_torneo, modalidad, puntos, habilitado)
			VALUES ($1, $2, $3, $4, 0, true)`

		_, err = tx.Exec(insertQuery, userID, equipo, torneoID, modalidad)
		if err != nil {
			return fmt.Errorf("error al inscribir usuario: %w", err)
		}

		// Actualizar estadísticas del usuario
		updateStatsQuery := `
			UPDATE user_stats
			SET torneos_participados = torneos_participados + 1
			WHERE user_id = $1`

		_, err = tx.Exec(updateStatsQuery, userID)
		if err != nil {
			return fmt.Errorf("error al actualizar estadísticas: %w", err)
		}

		return nil
	})
}
