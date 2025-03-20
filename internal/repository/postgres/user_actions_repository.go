package postgres

import (
	"backend_proyecto_verde/internal/models"
	"backend_proyecto_verde/pkg/database"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"

	"bytes"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type UserActionsRepository struct {
	db *sql.DB
}

func NewUserActionsRepository(db *sql.DB) *UserActionsRepository {
	return &UserActionsRepository{db: db}
}

func (r *UserActionsRepository) CreateAction(action *models.UserAction) error {
	return database.WithTransaction(r.db, func(tx *sql.Tx) error {
		query := `
			INSERT INTO user_actions (
				user_id, tipo_accion, foto, latitud, longitud, ciudad, lugar,
				en_colaboracion, colaboradores, es_para_torneo, id_torneo
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
			RETURNING id, created_at`

		err := tx.QueryRow(
			query,
			action.UserID,
			action.TipoAccion,
			action.Foto,
			action.Latitud,
			action.Longitud,
			action.Ciudad,
			action.Lugar,
			action.EnColaboracion,
			pq.Array(action.Colaboradores),
			action.EsParaTorneo,
			action.IDTorneo,
		).Scan(&action.ID, &action.CreatedAt)

		if err != nil {
			return err
		}

		// Actualizar estadísticas del usuario
		updateStatsQuery := `
			UPDATE user_stats
			SET acciones = acciones + 1, puntos = puntos + $1
			WHERE user_id = $2`

		var puntos int
		switch action.TipoAccion {
		case "ayuda":
			puntos = 50
		case "alerta":
			puntos = 40
		case "descubrimiento":
			puntos = 25
		default:
			puntos = 5
		}

		_, err = tx.Exec(updateStatsQuery, puntos, action.UserID)
		if err != nil {
			return err
		}

		// Si es para un torneo, actualizar puntos del torneo
		if action.EsParaTorneo && action.IDTorneo != nil {
			updateTorneoQuery := `
				UPDATE torneo_estadisticas
				SET puntos = puntos + $1
				WHERE id_jugador = $2 AND id_torneo = $3`

			_, err = tx.Exec(updateTorneoQuery, puntos, action.UserID, *action.IDTorneo)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *UserActionsRepository) GetTorneoID(userID string) (string, error) {
	var torneoID string
	query := `
		SELECT torneo_id
		FROM user_stats
		WHERE user_id = $1`

	err := r.db.QueryRow(query, userID).Scan(&torneoID)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("el usuario no tiene torneo asociado")
		}
		return "", fmt.Errorf("error al obtener el torneo del usuario: %w", err)
	}

	return torneoID, nil
}

func (r *UserActionsRepository) TorneoPoints(userID string, actionType string, torneoID string) (int, error) {
	var puntos int
	switch actionType {
	case "ayuda":
		puntos = 50
	case "alerta":
		puntos = 40
	case "descubrimiento":
		puntos = 25
	default:
		puntos = 5
	}

	// Actualizar puntos del torneo si es necesario
	updateTorneoQuery := `
		UPDATE torneo_estadisticas
		SET puntos = puntos + $1
		WHERE id_jugador = $2 AND id_torneo = $3`

	_, err := r.db.Exec(updateTorneoQuery, puntos, userID, torneoID)
	if err != nil {
		return 0, fmt.Errorf("error al actualizar puntos del torneo: %w", err)
	}

	return puntos, nil
}

func (r *UserActionsRepository) UploadImage(file multipart.File) (string, error) {
	if file == nil {
		return "", nil
	}
	defer file.Close()

	// Leer el contenido del archivo
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("error al leer la imagen: %v", err)
	}

	// Crear el body multipart
	boundary := "---011000010111000001101001"
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.SetBoundary(boundary)
	fileName := uuid.New().String() + ".jpg"

	part, err := writer.CreateFormFile("image", fileName)
	if err != nil {
		return "", fmt.Errorf("error al crear el formulario multipart: %v", err)
	}
	part.Write(fileBytes)
	writer.Close()

	// Crear la solicitud al CDN
	url := os.Getenv("CDN_URL") + "/api/cdn/upload/image/"
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return "", fmt.Errorf("error al crear la solicitud al CDN: %v", err)
	}

	// Agregar los headers requeridos
	req.Header.Set("Content-Type", "multipart/form-data; boundary="+boundary)
	req.Header.Set("Accept", "application/json")

	// Realizar la solicitud
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error al subir la imagen al CDN: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error en la respuesta del CDN: código %d", resp.StatusCode)
	}

	// Decodificar la respuesta para obtener la URL de la imagen
	var cdnResponse struct {
		URL string `json:"file_url"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&cdnResponse); err != nil {
		return "", fmt.Errorf("error al procesar la respuesta del CDN: %v", err)
	}

	imageURL := cdnResponse.URL
	if !strings.Contains(cdnResponse.URL, "/api/cdn") {
		imageURL = strings.Replace(cdnResponse.URL, strings.Replace(os.Getenv("CDN_URL"), "http://", "", 1), strings.Replace(os.Getenv("CDN_URL"), "http://", "", 1)+"/api/cdn", 1)
		fmt.Println("imageURL", imageURL)
	}

	return imageURL, nil
}

func (r *UserActionsRepository) GetUserActions(userID string) ([]models.UserAction, error) {
	query := `
		SELECT id, user_id, tipo_accion, foto, latitud, longitud, ciudad, lugar,
			en_colaboracion, colaboradores, es_para_torneo, id_torneo,
			created_at
		FROM user_actions
		WHERE user_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var actions []models.UserAction
	for rows.Next() {
		var a models.UserAction
		var colaboradores []string
		err := rows.Scan(
			&a.ID, &a.UserID, &a.TipoAccion, &a.Foto,
			&a.Latitud, &a.Longitud, &a.Ciudad, &a.Lugar,
			&a.EnColaboracion, pq.Array(&colaboradores), &a.EsParaTorneo, &a.IDTorneo,
			&a.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		if len(colaboradores) > 0 {
			a.Colaboradores = &colaboradores
		}
		actions = append(actions, a)
	}
	return actions, nil
}

func (r *UserActionsRepository) SoftDeleteAction(id string) error {
	return database.WithTransaction(r.db, func(tx *sql.Tx) error {
		// Obtener información de la acción antes de eliminarla
		var userID, tipoAccion string
		var esParaTorneo bool
		var idTorneo *string

		getActionQuery := `
			SELECT user_id, tipo_accion, es_para_torneo, id_torneo
			FROM user_actions
			WHERE id = $1 AND deleted_at IS NULL`

		err := tx.QueryRow(getActionQuery, id).Scan(&userID, &tipoAccion, &esParaTorneo, &idTorneo)
		if err != nil {
			return err
		}

		// Marcar la acción como eliminada
		updateQuery := `
			UPDATE user_actions
			SET deleted_at = $1
			WHERE id = $2 AND deleted_at IS NULL`

		result, err := tx.Exec(updateQuery, time.Now(), id)
		if err != nil {
			return err
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return err
		}

		if rowsAffected == 0 {
			return fmt.Errorf("no se encontró la acción o ya fue eliminada")
		}

		return nil
	})
}

func (r *UserActionsRepository) GetActionByID(id string) (*models.UserAction, error) {
	var action models.UserAction
	query := `SELECT id, user_id, tipo_accion, foto FROM user_actions WHERE id = $1 AND deleted_at IS NULL`
	err := r.db.QueryRow(query, id).Scan(&action.ID, &action.UserID, &action.TipoAccion, &action.Foto)
	if err != nil {
		return nil, err
	}
	return &action, nil
}

func (r *UserActionsRepository) DecrementUserPoints(userID string, tipoAccion string) error {
	var query string
	switch tipoAccion {
	case "ayuda":
		query = `UPDATE user_stats SET puntos = puntos - 50 WHERE user_id = $1`
	case "alerta":
		query = `UPDATE user_stats SET puntos = puntos - 40 WHERE user_id = $1`
	case "descubrimiento":
		query = `UPDATE user_stats SET puntos = puntos - 25 WHERE user_id = $1`
	default:
		return fmt.Errorf("tipo de acción no válido")
	}

	_, err := r.db.Exec(query, userID)
	return err
}

func (r *UserActionsRepository) DecrementTotalActions(userID string) error {
	query := `UPDATE user_stats SET acciones = acciones - 1 WHERE user_id = $1`
	_, err := r.db.Exec(query, userID)
	return err
}

func (r *UserActionsRepository) GetAllActions() ([]models.UserAction, error) {
	query := `
		SELECT id, user_id, tipo_accion, foto, latitud, longitud, ciudad, lugar,
			en_colaboracion, colaboradores, es_para_torneo, id_torneo,
			created_at
		FROM user_actions
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var actions []models.UserAction
	for rows.Next() {
		var a models.UserAction
		var colaboradores []string
		err := rows.Scan(
			&a.ID, &a.UserID, &a.TipoAccion, &a.Foto,
			&a.Latitud, &a.Longitud, &a.Ciudad, &a.Lugar,
			&a.EnColaboracion, pq.Array(&colaboradores), &a.EsParaTorneo, &a.IDTorneo,
			&a.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		if len(colaboradores) > 0 {
			a.Colaboradores = &colaboradores
		}
		actions = append(actions, a)
	}
	return actions, nil
}
