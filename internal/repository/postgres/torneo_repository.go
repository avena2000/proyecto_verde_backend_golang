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

func (r *TorneoRepository) GetIfTournamentOwner(userID string) (*models.UserStats, string, error) {
	// Obtener las estadísticas del usuario
	query := `
		SELECT id, user_id, puntos, acciones, torneos_participados, cantidad_amigos,
			es_dueno_torneo, torneos_ganados, pending_medalla, pending_amigo, torneo_id
		FROM user_stats
		WHERE user_id = $1`

	var stats models.UserStats
	err := r.db.QueryRow(query, userID).Scan(
		&stats.ID, &stats.UserID, &stats.Puntos,
		&stats.Acciones, &stats.TorneosParticipados, &stats.CantidadAmigos,
		&stats.EsDuenoTorneo, &stats.TorneosGanados, &stats.PendingMedalla,
		&stats.PendingAmigo, &stats.TorneoId,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			// Si no existen estadísticas, crear un objeto vacío
			stats = models.UserStats{
				UserID: userID,
			}
		} else {
			return nil, "Se encontró un problema, intenta nuevamente", fmt.Errorf("error al obtener estadísticas del usuario: %w", err)
		}
	}

	// Verificar si el usuario es dueño de algún torneo activo según la tabla de torneos
	queryTorneoActivo := `
		SELECT id, nombre
		FROM torneos
		WHERE id_creator = $1 AND finalizado = false`

	var torneoID, torneoNombre string
	err = r.db.QueryRow(queryTorneoActivo, userID).Scan(&torneoID, &torneoNombre)
	if err == nil {
		// Usuario es dueño de un torneo activo
		return &stats, "Eres dueño de un torneo activo: " + torneoNombre, fmt.Errorf("%s", torneoNombre)
	}

	// Verificar si el usuario está registrado en algún torneo
	queryRegistrado := `
		SELECT t.id, t.nombre
		FROM torneo_estadisticas te
		JOIN torneos t ON te.id_torneo = t.id
		WHERE te.id_jugador = $1 AND t.finalizado = false`

	err = r.db.QueryRow(queryRegistrado, userID).Scan(&torneoID, &torneoNombre)
	if err == nil {
		// Usuario está registrado en un torneo activo
		return &stats, "Ya estás registrado en un torneo activo: " + torneoNombre, fmt.Errorf("%s", torneoNombre)
	}

	// Verificar si en user_stats se marca que el usuario es dueño de un torneo
	if stats.EsDuenoTorneo {
		return &stats, "Según tus estadísticas, ya eres dueño de un torneo", fmt.Errorf("según tus estadísticas, ya eres dueño de un torneo")
	}

	// Si llegamos aquí, no hay problemas
	return &stats, "No hay problemas", nil
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

		// Actualizar estadísticas del usuario como dueño del torneo
		// Primero verificamos si existen estadísticas para este usuario
		checkQuery := `
			SELECT EXISTS (
				SELECT 1 FROM user_stats WHERE user_id = $1
			)`

		var exists bool
		err = tx.QueryRow(checkQuery, torneo.IDCreator).Scan(&exists)
		if err != nil {
			return fmt.Errorf("error al verificar estadísticas del usuario: %w", err)
		}

		if exists {
			// Actualizar estadísticas existentes
			updateQuery := `
				UPDATE user_stats
				SET es_dueno_torneo = true, torneo_id = $1
				WHERE user_id = $2`

			_, err = tx.Exec(updateQuery, torneo.ID, torneo.IDCreator)
			if err != nil {
				return fmt.Errorf("error al actualizar estadísticas del usuario: %w", err)
			}
		}

		return nil
	})
}

func (r *TorneoRepository) GetTorneoByCodeID(codeID string) (*models.Torneo, error) {
	query := `
		SELECT id, nombre, modalidad, ubicacion_a_latitud, ubicacion_a_longitud,
			nombre_ubicacion_a, ubicacion_b_latitud, ubicacion_b_longitud,
			nombre_ubicacion_b, fecha_inicio, fecha_fin, ubicacion_aproximada,
			metros_aproximados, finalizado, ganador_versus, ganador_individual
		FROM torneos
		WHERE code_id = $1`

	torneo := &models.Torneo{}
	err := r.db.QueryRow(query, codeID).Scan(
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
		WHERE id_creator = $1 AND finalizado = false`

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

			// Asegurarse de cerrar los rows inmediatamente después de su uso
			var equipoA, equipoB bool
			var puntosA, puntosB int

			for rows.Next() {
				var equipo bool
				var puntos int
				if err := rows.Scan(&equipo, &puntos); err != nil {
					rows.Close() // Cerrar explícitamente en caso de error
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

			// Cerrar filas y verificar errores
			err = rows.Err()
			rows.Close() // Cerrar explícitamente después de procesar
			if err != nil {
				return fmt.Errorf("error al procesar estadísticas: %w", err)
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

		// Actualizar estadísticas del usuario como dueño del torneo
		// Primero verificamos si existen estadísticas para este usuario
		checkQuery := `
			SELECT EXISTS (
				SELECT 1 FROM user_stats WHERE user_id = $1
			)`

		var exists bool
		err = tx.QueryRow(checkQuery, idCreator).Scan(&exists)
		if err != nil {
			return fmt.Errorf("error al verificar estadísticas del usuario: %w", err)
		}

		if exists {
			// Actualizar estadísticas existentes
			updateQuery := `
				UPDATE user_stats
				SET es_dueno_torneo = false, torneo_id = null
				WHERE user_id = $1`

			_, err = tx.Exec(updateQuery, idCreator)
			if err != nil {
				return fmt.Errorf("error al actualizar estadísticas del usuario: %w", err)
			}
		}

		//Actualizar estadísticas de torneos participados

		// Obtener todos los usuarios que participaron en el torneo
		participantesQuery := `
			SELECT id_jugador
			FROM torneo_estadisticas
			WHERE id_torneo = $1`

		participantesRows, err := tx.Query(participantesQuery, torneoID)
		if err != nil {
			return fmt.Errorf("error al obtener participantes del torneo: %w", err)
		}

		// Procesar los participantes y actualizar sus estadísticas
		var participanteIDs []string
		for participantesRows.Next() {
			var participanteID string
			if err := participantesRows.Scan(&participanteID); err != nil {
				participantesRows.Close() // Cerrar explícitamente en caso de error
				return fmt.Errorf("error al leer ID de participante: %w", err)
			}
			participanteIDs = append(participanteIDs, participanteID)
		}

		// Cerrar filas y verificar errores
		err = participantesRows.Err()
		participantesRows.Close() // Cerrar explícitamente después de procesar
		if err != nil {
			return fmt.Errorf("error al procesar participantes: %w", err)
		}

		// Actualizar estadísticas para cada participante
		for _, participanteID := range participanteIDs {
			updateParticipanteQuery := `
				UPDATE user_stats
				SET torneo_id = null
				WHERE user_id = $1`

			_, err = tx.Exec(updateParticipanteQuery, participanteID)
			if err != nil {
				return fmt.Errorf("error al actualizar estadísticas del participante %s: %w", participanteID, err)
			}
		}

		return nil
	})
}

func (r *TorneoRepository) BorrarTorneoEstadisticas(idCreator string) (string, error) {
	var message string
	err := database.WithTransaction(r.db, func(tx *sql.Tx) error {
		// Obtener el ID del torneo
		query := `
			SELECT id
			FROM torneos
			WHERE id_creator = $1 AND finalizado = false`

		var torneoID string
		err := tx.QueryRow(query, idCreator).Scan(&torneoID)
		if err != nil {
			if err == sql.ErrNoRows {
				message = "No se encontró un torneo activo para este usuario"
				return fmt.Errorf("no se encontró un torneo activo para este usuario")
			}
			message = "Error al buscar torneo"
			return fmt.Errorf("error al buscar torneo: %w", err)
		}

		// Verificar si hay registros en torneo_estadisticas para este torneo
		checkEstadisticasQuery := `
			SELECT COUNT(*)
			FROM torneo_estadisticas
			WHERE id_torneo = $1`

		var cantidadParticipantes int
		err = tx.QueryRow(checkEstadisticasQuery, torneoID).Scan(&cantidadParticipantes)
		if err != nil {
			message = "Error al verificar participantes"
			return fmt.Errorf("error al verificar participantes: %w", err)
		}

		// Si hay participantes, no permitir borrar el torneo
		if cantidadParticipantes > 0 {
			message = fmt.Sprintf("No se puede eliminar el torneo porque ya tiene %d participante(s)", cantidadParticipantes)
			return fmt.Errorf("no se puede eliminar el torneo porque ya tiene %d participante(s)", cantidadParticipantes)
		}

		// Eliminar estadísticas del torneo (por si acaso, aunque no debería haber)
		deleteQuery := `
			DELETE FROM torneo_estadisticas
			WHERE id_torneo = $1`

		_, err = tx.Exec(deleteQuery, torneoID)
		if err != nil {
			message = "Error al eliminar estadísticas"
			return fmt.Errorf("error al eliminar estadísticas: %w", err)
		}

		// Eliminar el torneo
		deleteTorneoQuery := `
			DELETE FROM torneos
			WHERE id = $1`

		_, err = tx.Exec(deleteTorneoQuery, torneoID)
		if err != nil {
			message = "Error al eliminar torneo"
			return fmt.Errorf("error al eliminar torneo: %w", err)
		}

		// Actualizar estadísticas del usuario creador
		updateStatsQuery := `
			UPDATE user_stats
			SET es_dueno_torneo = false, torneo_id = null
			WHERE user_id = $1`

		_, err = tx.Exec(updateStatsQuery, idCreator)
		if err != nil {
			message = "Error al actualizar estadísticas del creador"
			return fmt.Errorf("error al actualizar estadísticas del creador: %w", err)
		}

		return nil
	})
	return message, err
}

func (r *TorneoRepository) InscribirUsuario(codeID string, userID string, team bool) error {
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

		// Inscribir al usuario
		insertQuery := `
			INSERT INTO torneo_estadisticas (id_jugador, equipo, id_torneo, modalidad, puntos, habilitado)
			VALUES ($1, $2, $3, $4, 0, true)`

		_, err = tx.Exec(insertQuery, userID, team, torneoID, modalidad)
		if err != nil {
			return fmt.Errorf("error al inscribir usuario: %w", err)
		}

		// Actualizar estadísticas del usuario
		updateStatsQuery := `
			UPDATE user_stats
			SET torneos_participados = torneos_participados + 1, torneo_id = $1
			WHERE user_id = $2`

		_, err = tx.Exec(updateStatsQuery, torneoID, userID)
		if err != nil {
			return fmt.Errorf("error al actualizar estadísticas: %w", err)
		}

		return nil
	})
}

// GetTorneosRelacionadosUsuario obtiene todos los torneos relacionados con un usuario,
// tanto los que ha creado como en los que ha participado
func (r *TorneoRepository) GetTorneosRelacionadosUsuario(userID string) ([]models.TorneoResumen, error) {
	// Resultado final que contendrá todos los torneos relacionados con el usuario
	var torneos []models.TorneoResumen

	// Consulta para obtener torneos creados por el usuario
	queryCreados := `
		SELECT id, nombre
		FROM torneos
		WHERE id_creator = $1
		ORDER BY fecha_inicio DESC`

	rowsCreados, err := r.db.Query(queryCreados, userID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener torneos creados: %w", err)
	}
	defer rowsCreados.Close()

	// Procesar torneos creados
	for rowsCreados.Next() {
		var torneo models.TorneoResumen
		if err := rowsCreados.Scan(&torneo.ID, &torneo.Nombre); err != nil {
			return nil, fmt.Errorf("error al leer torneo creado: %w", err)
		}
		torneos = append(torneos, torneo)
	}

	// Consulta para obtener torneos en los que ha participado pero no ha creado
	queryParticipado := `
		SELECT t.id, t.nombre
FROM (
    SELECT DISTINCT ON (t.id) t.id, t.nombre, t.fecha_inicio
    FROM torneo_estadisticas te
    JOIN torneos t ON te.id_torneo = t.id
    WHERE te.id_jugador = $1 AND t.id_creator != $1
    ORDER BY t.id, t.fecha_inicio DESC
) AS t  -- Alias correcto para la subconsulta
ORDER BY t.fecha_inicio DESC;
`

	rowsParticipado, err := r.db.Query(queryParticipado, userID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener torneos participados: %w", err)
	}
	defer rowsParticipado.Close()

	// Procesar torneos participados
	for rowsParticipado.Next() {
		var torneo models.TorneoResumen
		if err := rowsParticipado.Scan(&torneo.ID, &torneo.Nombre); err != nil {
			return nil, fmt.Errorf("error al leer torneo participado: %w", err)
		}
		torneos = append(torneos, torneo)
	}

	if len(torneos) == 0 {
		return []models.TorneoResumen{}, nil
	}

	return torneos, nil
}

// SalirTorneo permite a un usuario abandonar un torneo en el que está participando
// eliminando sus estadísticas del torneo y actualizando sus estadísticas de usuario
func (r *TorneoRepository) SalirTorneo(userID string, torneoID string) error {
	return database.WithTransaction(r.db, func(tx *sql.Tx) error {
		// Verificar si el usuario está inscrito en el torneo
		checkQuery := `
			SELECT EXISTS (
				SELECT 1 FROM torneo_estadisticas
				WHERE id_torneo = $1 AND id_jugador = $2
			)`

		var exists bool
		err := tx.QueryRow(checkQuery, torneoID, userID).Scan(&exists)
		if err != nil {
			return fmt.Errorf("error al verificar inscripción: %w", err)
		}

		if !exists {
			return fmt.Errorf("el usuario no está inscrito en este torneo")
		}

		// Eliminar las estadísticas del usuario en este torneo
		deleteQuery := `
			DELETE FROM torneo_estadisticas
			WHERE id_torneo = $1 AND id_jugador = $2`

		_, err = tx.Exec(deleteQuery, torneoID, userID)
		if err != nil {
			return fmt.Errorf("error al eliminar estadísticas del torneo: %w", err)
		}

		// Actualizar estadísticas del usuario (reducir torneos_participados y quitar torneo_id)
		updateStatsQuery := `
			UPDATE user_stats
			SET torneos_participados = GREATEST(0, torneos_participados - 1), torneo_id = null
			WHERE user_id = $1`

		_, err = tx.Exec(updateStatsQuery, userID)
		if err != nil {
			return fmt.Errorf("error al actualizar estadísticas de usuario: %w", err)
		}

		return nil
	})
}

// GetEquipoUsuarioTorneo obtiene el equipo al que pertenece un usuario en un torneo específico
func (r *TorneoRepository) GetEquipoUsuarioTorneo(torneoID string, userID string) (*bool, error) {
	query := `
		SELECT equipo
		FROM torneo_estadisticas
		WHERE id_torneo = $1 AND id_jugador = $2`

	var equipo bool
	err := r.db.QueryRow(query, torneoID, userID).Scan(&equipo)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("el usuario no está inscrito en este torneo")
		}
		return nil, fmt.Errorf("error al obtener equipo del usuario: %w", err)
	}

	return &equipo, nil
}

// FindExpiredTournaments busca los torneos cuya fecha de fin ya ha pasado pero no están marcados como finalizados
func (r *TorneoRepository) FindExpiredTournaments() ([]string, error) {
	query := `
		SELECT id_creator
		FROM torneos
		WHERE fecha_fin < NOW() AND finalizado = false`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error al buscar torneos vencidos: %w", err)
	}
	defer rows.Close()

	var creatorIDs []string
	for rows.Next() {
		var creatorID string
		if err := rows.Scan(&creatorID); err != nil {
			return nil, fmt.Errorf("error al leer ID del creador: %w", err)
		}
		creatorIDs = append(creatorIDs, creatorID)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error al procesar resultados: %w", err)
	}

	return creatorIDs, nil
}

// UpdateTorneoFechaFin actualiza solo la fecha de fin de un torneo específico
func (r *TorneoRepository) UpdateTorneoFechaFin(id string, fechaFin string) error {
	query := `
		UPDATE torneos
		SET fecha_fin = $1
		WHERE id = $2`

	_, err := r.db.Exec(query, fechaFin, id)
	if err != nil {
		return fmt.Errorf("error al actualizar fecha de fin del torneo: %w", err)
	}

	return nil
}
