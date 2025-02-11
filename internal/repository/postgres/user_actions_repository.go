package postgres

import (
	"backend_proyecto_verde/internal/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
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
	query := `
		INSERT INTO user_actions (
			user_id, tipo_accion, foto, latitud, longitud,
			en_colaboracion, colaboradores, es_para_torneo, id_torneo
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at`

	err := r.db.QueryRow(
		query,
		action.UserID,
		action.TipoAccion,
		action.Foto,
		action.Latitud,
		action.Longitud,
		action.EnColaboracion,
		pq.Array(action.Colaboradores),
		action.EsParaTorneo,
		action.IDTorneo,
	).Scan(&action.ID, &action.CreatedAt)

	if err != nil {
		return err
	}

	puntos := 0

	if action.TipoAccion == "ayuda" {
		puntos = 50
	} else if action.TipoAccion == "alerta" {
		puntos = 40
	} else if action.TipoAccion == "descubrimiento" {
		puntos = 25
	}

	updateQuery := `
		UPDATE user_stats
		SET acciones = acciones + 1,
		puntos = puntos + $2
		WHERE user_id = $1;`

	_, err = r.db.Exec(updateQuery, action.UserID, puntos)
	if err != nil {
		return err
	}

	return nil

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
	url := "http://localhost:8080/api/cdn/upload/image/"
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
		imageURL = strings.Replace(cdnResponse.URL, "localhost:8080", "localhost:8080/api/cdn", 1)
	}

	return imageURL, nil
}

func (r *UserActionsRepository) GetUserActions(userID string, limit, offset int) ([]models.UserAction, error) {
	query := `
		SELECT id, user_id, tipo_accion, foto, latitud, longitud,
			en_colaboracion, colaboradores, es_para_torneo, id_torneo,
			created_at
		FROM user_actions
		WHERE user_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var actions []models.UserAction
	for rows.Next() {
		var a models.UserAction
		err := rows.Scan(
			&a.ID, &a.UserID, &a.TipoAccion, &a.Foto,
			&a.Latitud, &a.Longitud, &a.EnColaboracion,
			&a.Colaboradores, &a.EsParaTorneo, &a.IDTorneo,
			&a.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		actions = append(actions, a)
	}
	return actions, nil
}

func (r *UserActionsRepository) SoftDeleteAction(id string) error {
	query := `
		UPDATE user_actions
		SET deleted_at = $1
		WHERE id = $2 AND deleted_at IS NULL`

	result, err := r.db.Exec(query, time.Now(), id)
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
