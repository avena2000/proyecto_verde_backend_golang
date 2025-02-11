package postgres

import (
	"backend_proyecto_verde/internal/models"
	"database/sql"
	"time"
)

type MedallasRepository struct {
	db *sql.DB
	userStatsRepo *UserRepository
}

func NewMedallasRepository(db *sql.DB, userStatsRepo *UserRepository) *MedallasRepository {
	return &MedallasRepository{db: db, userStatsRepo: userStatsRepo}
}

func (r *MedallasRepository) CreateMedalla(medalla *models.Medalla) error {
	query := `
		INSERT INTO medallas (
			nombre, descripcion, dificultad, requiere_amistades,
			requiere_puntos, requiere_acciones, requiere_torneos,
			requiere_victoria_torneos, numero_requerido
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id`

	return r.db.QueryRow(
		query,
		medalla.Nombre,
		medalla.Descripcion,
		medalla.Dificultad,
		medalla.RequiereAmistades,
		medalla.RequierePuntos,
		medalla.RequiereAcciones,
		medalla.RequiereTorneos,
		medalla.RequiereVictoriaTorneos,
		medalla.NumeroRequerido,
	).Scan(&medalla.ID)
}

func (r *MedallasRepository) AutoAsignMedallas(userID string) (error) {
	// Obtener todas las medallas
	query := `SELECT id, nombre, descripcion, dificultad, requiere_amistades, requiere_puntos,
              requiere_acciones, requiere_torneos, requiere_victoria_torneos, numero_requerido
              FROM medallas`
	rows, err := r.db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	medallas := make([]models.Medalla, 0)
	for rows.Next() {
		var m models.Medalla
		if err := rows.Scan(&m.ID, &m.Nombre, &m.Descripcion, &m.Dificultad,
			&m.RequiereAmistades, &m.RequierePuntos, &m.RequiereAcciones,
			&m.RequiereTorneos, &m.RequiereVictoriaTorneos, &m.NumeroRequerido); err != nil {
			return err
		}
		medallas = append(medallas, m)
	}
	if err := rows.Err(); err != nil {
		return err
	}

	// Obtener medallas ya ganadas por el usuario
	queryAsignadas := `SELECT id_medalla FROM medallas_ganadas WHERE id_usuario = $1`
	rowsAsignadas, err := r.db.Query(queryAsignadas, userID)
	if err != nil {
		return err
	}
	defer rowsAsignadas.Close()

	medallasGanadas := make(map[string]bool)
	for rowsAsignadas.Next() {
		var idMedalla string
		if err := rowsAsignadas.Scan(&idMedalla); err != nil {
			return err
		}
		medallasGanadas[idMedalla] = true
	}

	// Filtrar medallas ya ganadas
	medallasPendientes := make([]models.Medalla, 0)
	for _, m := range medallas {
		if !medallasGanadas[m.ID] {
			medallasPendientes = append(medallasPendientes, m)
		}
	}

	// Obtener estadísticas del usuario
	userStats, err := r.userStatsRepo.GetUserStats(userID)
	if err != nil {
		return err
	}

	// Asignar medallas que cumplan con los requisitos
	medallasAsignadas := 0
	for _, m := range medallasPendientes {
		if cumpleRequisitos(m, *userStats) {
			if err := r.AsignarMedalla(userID, m.ID); err != nil {
				return err
			}
			medallasAsignadas++
		}
	}

	err = r.userStatsRepo.UpdateUserPendingMedalla(medallasAsignadas, userID)
	if err != nil {
		return err
	}

	return nil
}

// Función auxiliar para verificar si el usuario cumple con los requisitos de la medalla
func cumpleRequisitos(m models.Medalla, userStats models.UserStats) bool {
	if m.RequiereAmistades && userStats.CantidadAmigos >= *m.NumeroRequerido {
		return true
	}
	if m.RequierePuntos && userStats.Puntos >= *m.NumeroRequerido {
		return true
	}
	if m.RequiereAcciones && userStats.Acciones >= *m.NumeroRequerido {
		return true
	}
	if m.RequiereTorneos && userStats.TorneosParticipados >= *m.NumeroRequerido {
		return true
	}
	if m.RequiereVictoriaTorneos && userStats.TorneosGanados >= *m.NumeroRequerido {
		return true
	}
	return false
}


func (r *MedallasRepository) GetMedallas() ([]models.Medalla, error) {
	query := `SELECT * FROM medallas`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var medallas []models.Medalla
	for rows.Next() {
		var m models.Medalla
		err := rows.Scan(
			&m.ID, &m.Nombre, &m.Descripcion, &m.Dificultad,
			&m.RequiereAmistades, &m.RequierePuntos, &m.RequiereAcciones,
			&m.RequiereTorneos, &m.RequiereVictoriaTorneos, &m.NumeroRequerido,
		)
		if err != nil {
			return nil, err
		}
		medallas = append(medallas, m)
	}
	return medallas, nil
}

func (r *MedallasRepository) AsignarMedalla(userID, medallaID string) error {
	query := `
		INSERT INTO medallas_ganadas (id_usuario, id_medalla, fecha_ganada)
		VALUES ($1, $2, $3)`

	_, err := r.db.Exec(query, userID, medallaID, time.Now())
	return err
}

func (r *MedallasRepository) GetMedallasUsuario(userID string) ([]models.MedallaGanada, error) {
	query := `
		SELECT id, id_usuario, id_medalla, fecha_ganada
		FROM medallas_ganadas
		WHERE id_usuario = $1
		ORDER BY fecha_ganada DESC`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var medallas []models.MedallaGanada
	for rows.Next() {
		var m models.MedallaGanada
		err := rows.Scan(
			&m.ID, &m.IDUsuario, &m.IDMedalla, &m.FechaGanada,
		)
		if err != nil {
			return nil, err
		}
		medallas = append(medallas, m)
	}
	return medallas, nil
} 