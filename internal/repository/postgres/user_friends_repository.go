package postgres

import (
	"backend_proyecto_verde/internal/models"
	"backend_proyecto_verde/pkg/database"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrSelfFriendRequest   = errors.New("no te puedes agregar a ti mismo")
	ErrFriendRequestExists = errors.New("la solicitud de amistad ya existe")
)

type UserFriendsRepository struct {
	db *sql.DB
}

func NewUserFriendsRepository(db *sql.DB) *UserFriendsRepository {
	return &UserFriendsRepository{db: db}
}

func (r *UserFriendsRepository) SendFriendRequest(userID, friendIDRequest string) error {
	var friendID string
	query := `SELECT user_id FROM user_basic_info WHERE friend_id = $1`
	err := r.db.QueryRow(query, friendIDRequest).Scan(&friendID)
	if err != nil {
		if err == sql.ErrNoRows {
			return err
		}
		return err
	}

	if userID == friendID {
		return ErrSelfFriendRequest
	}

	orderedUserID := min(userID, friendID)
	orderedFriendID := max(userID, friendID)
	friendship, err := getFriendship(r.db, orderedUserID, orderedFriendID)
	if err != nil {
		return err
	}

	if friendship != "" {
		return ErrFriendRequestExists
	}

	return database.WithTransaction(r.db, func(tx *sql.Tx) error {
		query := `
			INSERT INTO user_friends (user_id, friend_id, pending_id)
			VALUES ($1, $2, $3)
			ON CONFLICT (user_id, friend_id) DO UPDATE
			SET deleted_at = NULL, pending_id = $3
			WHERE user_friends.deleted_at IS NOT NULL`

		_, err := tx.Exec(query, orderedUserID, orderedFriendID, friendID)
		if err != nil {
			return err
		}

		query = `
			UPDATE user_stats
			SET pending_amigo = pending_amigo + 1
			WHERE user_id = $1
		`
		_, err = tx.Exec(query, friendID)
		if err != nil {
			return err
		}

		return nil
	})
}

func min(a, b string) string {
	if a < b {
		return a
	}
	return b
}

func max(a, b string) string {
	if a > b {
		return a
	}
	return b
}

func getFriendship(db *sql.DB, userID, friendID string) (string, error) {
	query := `
		SELECT
		    CASE
		        WHEN user_id = $1 THEN friend_id
		        ELSE user_id
		    END AS friend
		FROM user_friends
		WHERE ((user_id = $1 AND friend_id = $2) OR (user_id = $2 AND friend_id = $1))
		AND deleted_at IS NULL
		LIMIT 1;
	`
	var friend string
	err := db.QueryRow(query, userID, friendID).Scan(&friend)
	if err == sql.ErrNoRows {
		return "", nil // No existe amistad
	} else if err != nil {
		return "", err // Otro error en la consulta
	}

	return friend, nil
}

func (r *UserFriendsRepository) AcceptFriendRequest(userID, friendID string) error {
	return database.WithTransaction(r.db, func(tx *sql.Tx) error {
		query := `
			UPDATE user_friends
			SET pending_id = NULL
			WHERE
			((user_id = $1 AND friend_id = $2) OR (user_id = $2 AND friend_id = $1))
			AND (pending_id = $1 OR pending_id = $2);
		`

		result, err := tx.Exec(query, userID, friendID)
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

		query = `
			UPDATE user_stats
			SET cantidad_amigos = cantidad_amigos + 1
			WHERE user_id IN ($1, $2);
		`
		_, err = tx.Exec(query, userID, friendID)
		if err != nil {
			return err
		}

		query = `
			UPDATE user_stats
			SET pending_amigo = pending_amigo - 1
			WHERE user_id = $1
		`
		_, err = tx.Exec(query, userID)
		if err != nil {
			return err
		}

		return nil
	})
}

func (r *UserFriendsRepository) GetFriendsList(userID string) ([]models.UserFriend, error) {
	query := `
		SELECT
			uf.id,
			CASE
				WHEN uf.user_id = $1 THEN uf.friend_id
				ELSE uf.user_id
			END AS friend_id,
			ubi.nombre,
			ubi.apellido,
			uf.pending_id,
			up.slogan,
			up.cabello,
			up.vestimenta,
			up.barba,
			up.detalle_facial,
			up.detalle_adicional
		FROM user_friends uf
		JOIN user_basic_info ubi
			ON ubi.user_id = CASE
								WHEN uf.user_id = $1 THEN uf.friend_id
								ELSE uf.user_id
							END
		JOIN user_profile up
			ON up.user_id = CASE
								WHEN uf.user_id = $1 THEN uf.friend_id
								ELSE uf.user_id
							END
		WHERE (uf.user_id = $1 OR uf.friend_id = $1)
		AND uf.deleted_at IS NULL;
`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	friends := []models.UserFriend{}
	for rows.Next() {
		var friend models.UserFriend
		err := rows.Scan(&friend.ID, &friend.FriendID, &friend.Nombre, &friend.Apellido, &friend.PendingID, &friend.Slogan, &friend.Cabello, &friend.Vestimenta, &friend.Barba, &friend.DetalleFacial, &friend.DetalleAdicional)
		if err != nil {
			return nil, err
		}
		friends = append(friends, friend)
	}

	return friends, nil
}

func (r *UserFriendsRepository) RemoveFriend(userID, friendID string) error {
	// Primero verificamos si la amistad existe y no ha sido borrada
	var existsAndNotDeleted bool
	searchQuery := `
		SELECT EXISTS (
			SELECT 1 FROM user_friends
			WHERE ((user_id = $1 AND friend_id = $2) OR (user_id = $2 AND friend_id = $1))
			AND deleted_at IS NULL
		)`
	err := r.db.QueryRow(searchQuery, userID, friendID).Scan(&existsAndNotDeleted)
	if err != nil {
		return err
	}

	// Si no existe o ya fue borrada, no hacemos nada
	if !existsAndNotDeleted {
		return nil
	}

	// Verificamos si es una amistad confirmada o pendiente
	var isPendingFriendship bool
	pendingQuery := `
		SELECT EXISTS (
			SELECT 1 FROM user_friends
			WHERE ((user_id = $1 AND friend_id = $2) OR (user_id = $2 AND friend_id = $1))
			AND deleted_at IS NULL
			AND pending_id IS NOT NULL
		)`
	err = r.db.QueryRow(pendingQuery, userID, friendID).Scan(&isPendingFriendship)
	if err != nil {
		return err
	}

	return database.WithTransaction(r.db, func(tx *sql.Tx) error {
		query := `
			UPDATE user_friends
			SET deleted_at = $1
			WHERE
			((user_id = $2 AND friend_id = $3) OR (user_id = $3 AND friend_id = $2))
			AND deleted_at IS NULL;
		`

		result, err := tx.Exec(query, time.Now(), userID, friendID)
		if err != nil {
			return err
		}

		// Verificamos que realmente se actualizó alguna fila
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return err
		}

		// Si no se actualizó ninguna fila, no actualizamos los contadores
		if rowsAffected == 0 {
			return nil
		}

		if !isPendingFriendship {
			// Es una amistad confirmada
			query = `
			UPDATE user_stats
			SET cantidad_amigos = cantidad_amigos - 1
			WHERE user_id IN ($1, $2);
			`
			_, err = tx.Exec(query, userID, friendID)
		} else {
			// Es una solicitud pendiente
			// Determinamos quién es el destinatario de la solicitud
			var pendingUserID string
			pendingQuery := `
				SELECT pending_id 
				FROM user_friends
				WHERE ((user_id = $1 AND friend_id = $2) OR (user_id = $2 AND friend_id = $1))
				AND deleted_at IS NULL
			`
			err = tx.QueryRow(pendingQuery, userID, friendID).Scan(&pendingUserID)
			if err != nil {
				return err
			}

			query = `
			UPDATE user_stats
			SET pending_amigo = pending_amigo - 1
			WHERE user_id = $1
			`
			_, err = tx.Exec(query, pendingUserID)
		}

		return err
	})
}
