package postgres

import (
	"backend_proyecto_verde/internal/models"
	"backend_proyecto_verde/internal/utils"
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
	query := `
		INSERT INTO user_access (username, password)
		VALUES ($1, $2)
		RETURNING id, username`

	var createdUser models.UserAccess
	err := r.db.QueryRow(query, user.Username, user.Password).
		Scan(&createdUser.ID, &createdUser.Username)

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
	query := `
		UPDATE user_access
		SET username = $1
		WHERE id = $2`

	_, err := r.db.Exec(query, user.Username, user.ID)
	return err
}

func (r *UserRepository) CreateUserBasicInfo(user *models.UserBasicInfo) error {
	query := `
		INSERT INTO user_basic_info (user_id, numero, nombre, apellido)
		VALUES ($1, $2, $3, $4)
		RETURNING id`

	_, err := r.db.Exec(query, user.ID, user.Numero, user.Nombre, user.Apellido)
	return err

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
	var exists bool

	checkQuery := `
		SELECT EXISTS (
			SELECT 1 FROM user_basic_info WHERE user_id = $1
		)`
	err := r.db.QueryRow(checkQuery, user.UserID).Scan(&exists)
	if err != nil {
		return err
	}

	if !exists {
		friendId := utils.GenerateUniqueFriendId(r.db, false)
		user.FriendId = friendId
	}

	query := `
		INSERT INTO user_basic_info (user_id, numero, nombre, apellido, friend_id)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (user_id) DO UPDATE
		SET nombre = $3, apellido = $4`

	_, err = r.db.Exec(query,
		user.UserID, user.Numero, user.Nombre, user.Apellido, user.FriendId,
	)
	return err

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

func (r *UserRepository) UpdateUserProfile(profile *models.UserProfile) error {

	query := `
		INSERT INTO user_profile (user_id, slogan, cabello, vestimenta, barba, detalle_facial, detalle_adicional)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (user_id) DO UPDATE
		SET slogan = $2, cabello = $3, vestimenta = $4, barba = $5, detalle_facial = $6, detalle_adicional = $7`

	_, err := r.db.Exec(query,
		profile.UserID, profile.Slogan, profile.Cabello, profile.Vestimenta,
		profile.Barba, profile.DetalleFacial, profile.DetalleAdicional,
	)
	return err

}

func (r *UserRepository) EditUserProfile(profile *models.EditProfile) error {

	query := `
		UPDATE user_profile
		SET slogan = $2
		WHERE user_id = $1`

	_, err := r.db.Exec(query, profile.UserID, profile.Slogan)

	if err != nil {
		return err
	}

	query = `
		UPDATE user_basic_info

		SET nombre = $2, apellido = $3
		WHERE user_id = $1`

	_, err = r.db.Exec(query, profile.UserID, profile.Nombre, profile.Apellido)

	return err
}

func (r *UserRepository) GetUserStats(userID string) (*models.UserStats, error) {
	query := `
		SELECT id, user_id, puntos, acciones,torneos_participados, cantidad_amigos,
			es_dueno_torneo, torneos_ganados, pending_medalla, pending_amigo
		FROM user_stats
		WHERE user_id = $1`

	var stats models.UserStats
	err := r.db.QueryRow(query, userID).Scan(
		&stats.ID, &stats.UserID, &stats.Puntos,
		&stats.Acciones, &stats.TorneosParticipados, &stats.CantidadAmigos,
		&stats.EsDuenoTorneo, &stats.TorneosGanados, &stats.PendingMedalla,
		&stats.PendingAmigo,
	)

	if err != nil {
		return nil, err
	}

	return &stats, nil
}

func (r *UserRepository) UpdateUserStats(stats *models.UserStats) error {
	query := `
		INSERT INTO user_stats (
			user_id, puntos, acciones, torneos_participados, cantidad_amigos,
			es_dueno_torneo
		) VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (user_id) DO UPDATE
		SET puntos = $2, acciones = $3, torneos_participados = $4, cantidad_amigos = $5,
			es_dueno_torneo = $6
	`

	_, err := r.db.Exec(query,
		stats.UserID, stats.Puntos, stats.Acciones, stats.TorneosParticipados,
		stats.CantidadAmigos, stats.EsDuenoTorneo,
	)
	if err != nil {
		fmt.Println("SQL Execution Error:", err)
	}
	return err
}

func (r *UserRepository) UpdateUserPendingMedalla(pending int, userID string) error {
	query := `
		INSERT INTO user_stats (
			user_id, pending_medalla
		) VALUES ($1, $2)
		ON CONFLICT (user_id) DO UPDATE
		SET pending_medalla = user_stats.pending_medalla + $2`

	_, err := r.db.Exec(query, userID, pending)
	if err != nil {
		fmt.Println("SQL Execution Error:", err)
	}
	return err
}

func (r *UserRepository) UpdateUserPendingAmigo(pending int, userID string) error {
	query := `
		INSERT INTO user_stats (
			user_id, pending_amigo
		) VALUES ($1, $2)
		ON CONFLICT (user_id) DO UPDATE
		SET pending_amigo = $2`

	_, err := r.db.Exec(query,
		userID, pending,
	)
	return err
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
		LEFT JOIN user_profile up ON us.user_id = up.user_id
		LEFT JOIN user_basic_info ub ON us.user_id = ub.user_id
		ORDER BY us.puntos DESC`

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
