package postgres

import (
	"backend_proyecto_verde/internal/models"
	"backend_proyecto_verde/internal/utils"
	"backend_proyecto_verde/pkg/database"
	"database/sql"
	"errors"
	"fmt"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(user *models.CreateUserAccess) (*models.UserAccess, error) {
	var createdUser models.UserAccess

	err := database.WithTransaction(r.db, func(tx *sql.Tx) error {
		query := `
			INSERT INTO user_access (username, password)
			VALUES ($1, $2)
			RETURNING id, username`

		err := tx.QueryRow(query, user.Username, user.Password).
			Scan(&createdUser.ID, &createdUser.Username)

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &createdUser, nil
}

func (r *UserRepository) ListUsers(limit, offset int) ([]models.UserAccess, error) {
	query := `
		SELECT id, username
		FROM user_access
		LIMIT $1 OFFSET $2`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.UserAccess
	for rows.Next() {
		var u models.UserAccess
		if err := rows.Scan(&u.ID, &u.Username); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *UserRepository) GetUserByID(id string) (*models.UserAccess, error) {
	query := `
		SELECT id, username
		FROM user_access
		WHERE id = $1`

	var user models.UserAccess
	err := r.db.QueryRow(query, id).Scan(&user.ID, &user.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("usuario no encontrado")
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) ReLoginUserByID(id string) (*models.LoginUserAccess, error) {
	query := `
	SELECT id, username
	FROM user_access
	WHERE id = $1`

	var user models.UserAccess
	err := r.db.QueryRow(query, id).Scan(&user.ID, &user.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("usuario no encontrado")
		}
		return nil, err
	}

	query = `
	SELECT EXISTS (
		SELECT 1 FROM user_basic_info WHERE user_id = $1
	)`

	userChecked := models.LoginUserAccess{
		ID:                    user.ID,
		Username:              user.Username,
		IsPersonalInformation: false,
	}

	var exists bool
	err = r.db.QueryRow(query, user.ID).Scan(&exists)
	if err == nil && exists {
		userChecked.IsPersonalInformation = true
	}

	return &userChecked, nil
}

func (r *UserRepository) GetUserByUsername(username string) (*models.UserAccess, error) {
	query := `
		SELECT id, username
		FROM user_access
		WHERE username = $1`

	var user models.UserAccess
	err := r.db.QueryRow(query, username).Scan(&user.ID, &user.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("usuario no encontrado")
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetUserByUsernameAndPassword(username string, password string) (*models.LoginUserAccess, error) {
	query := `
		SELECT id, username
		FROM user_access
		WHERE username = $1 AND password = $2`

	var user models.UserAccess
	err := r.db.QueryRow(query, username, password).Scan(&user.ID, &user.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("usuario no encontrado")
		}
		return nil, err
	}

	query = `
	SELECT EXISTS (
		SELECT 1 FROM user_basic_info WHERE user_id = $1
	)`

	userChecked := models.LoginUserAccess{
		ID:                    user.ID,
		Username:              user.Username,
		IsPersonalInformation: false,
	}

	var exists bool
	err = r.db.QueryRow(query, user.ID).Scan(&exists)
	if err == nil && exists {
		userChecked.IsPersonalInformation = true
	}

	return &userChecked, nil
}

func (r *UserRepository) UpdateUser(user *models.UserAccess) error {
	return database.WithTransaction(r.db, func(tx *sql.Tx) error {
		query := `
			UPDATE user_access
			SET username = $1
			WHERE id = $2`

		_, err := tx.Exec(query, user.Username, user.ID)
		return err
	})
}

func (r *UserRepository) CreateUserBasicInfo(user *models.UserBasicInfo) error {
	return database.WithTransaction(r.db, func(tx *sql.Tx) error {
		query := `
			INSERT INTO user_basic_info (user_id, numero, nombre, apellido)
			VALUES ($1, $2, $3, $4)
			RETURNING id`

		_, err := tx.Exec(query, user.ID, user.Numero, user.Nombre, user.Apellido)
		return err
	})
}

func (r *UserRepository) GetUserBasicInfo(user *models.UserBasicInfo) (*models.UserBasicInfo, error) {
	query := `
		SELECT id, user_id, numero, nombre, apellido, friend_id
		FROM user_basic_info
		WHERE user_id = $1`

	var userBasicInfo models.UserBasicInfo
	err := r.db.QueryRow(query, user.UserID).Scan(&userBasicInfo.ID, &userBasicInfo.UserID, &userBasicInfo.Numero, &userBasicInfo.Nombre, &userBasicInfo.Apellido, &userBasicInfo.FriendId)
	if err != nil {
		return nil, err
	}

	return &userBasicInfo, nil
}

func (r *UserRepository) CreateOrUpdateUserBasicInfo(user *models.UserBasicInfo) error {
	return database.WithTransaction(r.db, func(tx *sql.Tx) error {
		// Verificar si ya existe información básica para este usuario
		query := `
			SELECT EXISTS (
				SELECT 1 FROM user_basic_info WHERE user_id = $1
			)`

		var exists bool
		err := tx.QueryRow(query, user.UserID).Scan(&exists)
		if err != nil {
			return err
		}

		if !exists {
			friendId := utils.GenerateUniqueFriendId(r.db, false)
			user.FriendId = friendId
		}

		if exists {
			// Actualizar información existente
			query = `
				UPDATE user_basic_info
				SET numero = $1, nombre = $2, apellido = $3
				WHERE user_id = $4`

			_, err = tx.Exec(query, user.Numero, user.Nombre, user.Apellido, user.UserID)
			return err
		} else {
			// Crear nueva información
			query = `
				INSERT INTO user_basic_info (user_id, numero, nombre, apellido, friend_id)
				VALUES ($1, $2, $3, $4, $5)`

			_, err = tx.Exec(query, user.UserID, user.Numero, user.Nombre, user.Apellido, user.FriendId)
			return err
		}
	})
}

func (r *UserRepository) GetUserProfile(userID string) (*models.UserProfile, error) {
	query := `
		SELECT id, user_id, slogan, cabello, vestimenta, barba, detalle_facial, detalle_adicional
		FROM user_profile
		WHERE user_id = $1`

	var profile models.UserProfile
	err := r.db.QueryRow(query, userID).Scan(
		&profile.ID, &profile.UserID, &profile.Slogan,
		&profile.Cabello, &profile.Vestimenta, &profile.Barba,
		&profile.DetalleFacial, &profile.DetalleAdicional,
	)
	if err != nil {
		return nil, err
	}

	return &profile, nil
}

func (r *UserRepository) CreateUserProfile(profile *models.UserProfile) error {
	return database.WithTransaction(r.db, func(tx *sql.Tx) error {
		// Verificar si ya existe un perfil para este usuario
		query := `
			SELECT EXISTS (
				SELECT 1 FROM user_profile WHERE user_id = $1
			)`

		var exists bool
		err := tx.QueryRow(query, profile.UserID).Scan(&exists)
		if err != nil {
			return err
		}

		if exists {
			// Actualizar perfil existente
			query = `
				UPDATE user_profile
				SET slogan = $1, cabello = $2, vestimenta = $3, barba = $4, detalle_facial = $5, detalle_adicional = $6
				WHERE user_id = $7`

			_, err = tx.Exec(query, profile.Slogan, profile.Cabello, profile.Vestimenta, profile.Barba, profile.DetalleFacial, profile.DetalleAdicional, profile.UserID)
			return err
		} else {
			// Crear nuevo perfil
			query = `
				INSERT INTO user_profile (user_id, slogan, cabello, vestimenta, barba, detalle_facial, detalle_adicional)
				VALUES ($1, $2, $3, $4, $5, $6, $7)`

			_, err = tx.Exec(query, profile.UserID, profile.Slogan, profile.Cabello, profile.Vestimenta, profile.Barba, profile.DetalleFacial, profile.DetalleAdicional)
			return err
		}
	})
}

func (r *UserRepository) EditUserProfile(profile *models.EditProfile) error {
	return database.WithTransaction(r.db, func(tx *sql.Tx) error {
		// Verificar si ya existe un perfil para este usuario
		query := `
			SELECT EXISTS (
				SELECT 1 FROM user_profile WHERE user_id = $1
			)`

		var exists bool
		err := tx.QueryRow(query, profile.UserID).Scan(&exists)
		if err != nil {
			return err
		}

		if !exists {
			return errors.New("el perfil no existe")
		}

		// Construir la consulta dinámicamente basada en los campos proporcionados
		query = "UPDATE user_profile SET "
		params := []interface{}{}
		paramCount := 1

		if profile.Slogan != nil {
			query += fmt.Sprintf("slogan = $%d, ", paramCount)
			params = append(params, *profile.Slogan)
			paramCount++
		}

		if profile.Cabello != nil {
			query += fmt.Sprintf("cabello = $%d, ", paramCount)
			params = append(params, *profile.Cabello)
			paramCount++
		}

		if profile.Vestimenta != nil {
			query += fmt.Sprintf("vestimenta = $%d, ", paramCount)
			params = append(params, *profile.Vestimenta)
			paramCount++
		}

		if profile.Barba != nil {
			query += fmt.Sprintf("barba = $%d, ", paramCount)
			params = append(params, *profile.Barba)
			paramCount++
		}

		if profile.DetalleFacial != nil {
			query += fmt.Sprintf("detalle_facial = $%d, ", paramCount)
			params = append(params, *profile.DetalleFacial)
			paramCount++
		}

		if profile.DetalleAdicional != nil {
			query += fmt.Sprintf("detalle_adicional = $%d, ", paramCount)
			params = append(params, *profile.DetalleAdicional)
			paramCount++
		}

		// Eliminar la última coma y espacio
		query = query[:len(query)-2]

		// Añadir la condición WHERE
		query += fmt.Sprintf(" WHERE user_id = $%d", paramCount)
		params = append(params, profile.UserID)

		_, err = tx.Exec(query, params...)
		return err
	})
}

func (r *UserRepository) GetUserStats(userID string) (*models.UserStats, error) {
	query := `
		SELECT id, user_id, puntos, acciones,torneos_participados, cantidad_amigos,
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
		return nil, err
	}

	return &stats, nil
}

func (r *UserRepository) UpdateUserStats(stats *models.UserStats) error {
	return database.WithTransaction(r.db, func(tx *sql.Tx) error {
		// Verificar si ya existen estadísticas para este usuario
		query := `
			SELECT EXISTS (
				SELECT 1 FROM user_stats WHERE user_id = $1
			)`

		var exists bool
		err := tx.QueryRow(query, stats.UserID).Scan(&exists)
		if err != nil {
			return err
		}

		if exists {
			// Actualizar estadísticas existentes
			query = `
				UPDATE user_stats
				SET puntos = $1, acciones = $2, torneos_participados = $3, cantidad_amigos = $4, 
					es_dueno_torneo = $5, torneos_ganados = $6, pending_medalla = $7, pending_amigo = $8
				WHERE user_id = $9`

			_, err = tx.Exec(query, stats.Puntos, stats.Acciones, stats.TorneosParticipados,
				stats.CantidadAmigos, stats.EsDuenoTorneo, stats.TorneosGanados,
				stats.PendingMedalla, stats.PendingAmigo, stats.UserID)
			return err
		} else {
			// Crear nuevas estadísticas
			query = `
				INSERT INTO user_stats (user_id, puntos, acciones, torneos_participados, cantidad_amigos, 
					es_dueno_torneo, torneos_ganados, pending_medalla, pending_amigo)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

			_, err = tx.Exec(query, stats.UserID, stats.Puntos, stats.Acciones, stats.TorneosParticipados,
				stats.CantidadAmigos, stats.EsDuenoTorneo, stats.TorneosGanados,
				stats.PendingMedalla, stats.PendingAmigo)
			return err
		}
	})
}

func (r *UserRepository) UpdateUserPendingMedalla(pending int, userID string) error {
	return database.WithTransaction(r.db, func(tx *sql.Tx) error {
		query := `
			UPDATE user_stats
			SET pending_medalla = $1
			WHERE user_id = $2`

		_, err := tx.Exec(query, pending, userID)
		return err
	})
}

func (r *UserRepository) UpdateUserPendingAmigo(pending int, userID string) error {
	return database.WithTransaction(r.db, func(tx *sql.Tx) error {
		query := `
			UPDATE user_stats
			SET pending_amigo = $1
			WHERE user_id = $2`

		_, err := tx.Exec(query, pending, userID)
		return err
	})
}

func (r *UserRepository) GetRanking() ([]models.UserRanking, error) {
	query := `
		SELECT
			us.user_id,
			us.puntos,
			us.acciones,
			us.torneos_ganados,
			us.cantidad_amigos,
			up.slogan,
			up.cabello,
			up.vestimenta,
			up.barba,
			up.detalle_facial,
			up.detalle_adicional,
			ub.nombre,
			ub.apellido
		FROM user_stats us
		LEFT JOIN user_profile up
		ON us.user_id = up.user_id
		LEFT JOIN user_basic_info ub
		ON us.user_id = ub.user_id
		WHERE ub.nombre IS NOT NULL
		AND ub.nombre <> ''
		AND ub.apellido IS NOT NULL
		AND ub.apellido <> ''
		AND up.slogan IS NOT NULL
		AND up.slogan <> ''
		AND up.cabello IS NOT NULL
		AND up.cabello <> ''
		AND up.vestimenta IS NOT NULL
		AND up.vestimenta <> ''
		AND up.barba IS NOT NULL
		AND up.barba <> ''
		AND up.detalle_facial IS NOT NULL
		AND up.detalle_facial <> ''
		AND up.detalle_adicional IS NOT NULL
		AND up.detalle_adicional <> ''
		ORDER BY us.puntos DESC;
`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ranking []models.UserRanking
	for rows.Next() {
		var r models.UserRanking
		err := rows.Scan(
			&r.UserID,
			&r.Puntos,
			&r.Acciones,
			&r.TorneosGanados,
			&r.CantidadAmigos,
			&r.Slogan,
			&r.Cabello,
			&r.Vestimenta,
			&r.Barba,
			&r.DetalleFacial,
			&r.DetalleAdicional,
			&r.Nombre,
			&r.Apellido,
		)
		if err != nil {
			return nil, err
		}
		ranking = append(ranking, r)
	}

	return ranking, nil
}

// GetRankingTorneo obtiene el ranking de usuarios para un torneo específico
// ordenado por los puntos obtenidos en ese torneo
func (r *UserRepository) GetRankingTorneo(torneoID string) ([]models.UserRanking, error) {
	query := `
		SELECT
			te.id_jugador AS user_id,
			te.puntos,
			us.acciones,
			us.torneos_ganados,
			us.cantidad_amigos,
			up.slogan,
			up.cabello,
			up.vestimenta,
			up.barba,
			up.detalle_facial,
			up.detalle_adicional,
			ub.nombre,
			ub.apellido
		FROM torneo_estadisticas te
		JOIN user_stats us ON te.id_jugador = us.user_id
		LEFT JOIN user_profile up ON te.id_jugador = up.user_id
		LEFT JOIN user_basic_info ub ON te.id_jugador = ub.user_id
		WHERE te.id_torneo = $1
		AND ub.nombre IS NOT NULL
		AND ub.nombre <> ''
		AND ub.apellido IS NOT NULL
		AND ub.apellido <> ''
		AND up.slogan IS NOT NULL
		AND up.slogan <> ''
		AND up.cabello IS NOT NULL
		AND up.cabello <> ''
		AND up.vestimenta IS NOT NULL
		AND up.vestimenta <> ''
		AND up.barba IS NOT NULL
		AND up.barba <> ''
		AND up.detalle_facial IS NOT NULL
		AND up.detalle_facial <> ''
		AND up.detalle_adicional IS NOT NULL
		AND up.detalle_adicional <> ''
		ORDER BY te.puntos DESC;
	`
	rows, err := r.db.Query(query, torneoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ranking []models.UserRanking
	for rows.Next() {
		var r models.UserRanking
		err := rows.Scan(
			&r.UserID,
			&r.Puntos,
			&r.Acciones,
			&r.TorneosGanados,
			&r.CantidadAmigos,
			&r.Slogan,
			&r.Cabello,
			&r.Vestimenta,
			&r.Barba,
			&r.DetalleFacial,
			&r.DetalleAdicional,
			&r.Nombre,
			&r.Apellido,
		)
		if err != nil {
			return nil, err
		}
		ranking = append(ranking, r)
	}

	return ranking, nil
}
