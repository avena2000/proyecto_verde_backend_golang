package postgres

import (
	"backend_proyecto_verde/internal/models"
	"database/sql"
	"time"
)

type UserFriendsRepository struct {
	db *sql.DB
}

func NewUserFriendsRepository(db *sql.DB) *UserFriendsRepository {
	return &UserFriendsRepository{db: db}
}

func (r *UserFriendsRepository) SendFriendRequest(userID, friendID string) error {
	query := `
		INSERT INTO user_friends (user_id, friend_id, pending_id)
		VALUES ($1, $2, $2)
		ON CONFLICT (user_id, friend_id) DO UPDATE
		SET deleted_at = NULL, pending_id = $2
		WHERE user_friends.deleted_at IS NOT NULL`

	_, err := r.db.Exec(query, userID, friendID)
	return err
}

func (r *UserFriendsRepository) AcceptFriendRequest(userID, friendID string) error {
	query := `
		UPDATE user_friends
		SET pending_id = NULL
		WHERE user_id = $1 AND friend_id = $2 AND pending_id = $1`

	result, err := r.db.Exec(query, userID, friendID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *UserFriendsRepository) GetFriendsList(userID string) ([]models.UserBasicInfo, error) {
	query := `
		SELECT u.id, u.nombre, u.apellido
		FROM user_basic_info u
		INNER JOIN user_friends f ON u.user_id = f.friend_id
		WHERE f.user_id = $1 AND f.pending_id IS NULL AND f.deleted_at IS NULL`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var friends []models.UserBasicInfo
	for rows.Next() {
		var f models.UserBasicInfo
		err := rows.Scan(&f.ID, &f.Nombre, &f.Apellido)
		if err != nil {
			return nil, err
		}
		friends = append(friends, f)
	}
	return friends, nil
}

func (r *UserFriendsRepository) RemoveFriend(userID, friendID string) error {
	query := `
		UPDATE user_friends
		SET deleted_at = $1
		WHERE user_id = $2 AND friend_id = $3 AND deleted_at IS NULL`

	_, err := r.db.Exec(query, time.Now(), userID, friendID)
	return err
} 