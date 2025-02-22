package utils

import (
	"database/sql"
	"log"
	"math/rand"
	"time"
)



func GenerateUniqueFriendId(db *sql.DB, tornament bool) string {
	var existingSequences = make(map[string]bool) // Almacena secuencias para evitar duplicados
	existingSequences = getExistingSequences(db, tornament)
	return generateUniqueSequence(5, existingSequences)
}

// Genera una secuencia aleatoria de movimientos de longitud `n`
func generateUniqueSequence(n int, existingSequences map[string]bool) string {
	moves := []string{"up", "down", "right", "left"}

	rand.Seed(time.Now().UnixNano())

	for {
		seq := ""
		for i := 0; i < n; i++ {
			seq += moves[rand.Intn(len(moves))] + "-"
		}
		seq = seq[:len(seq)-1] // Remueve el último "-"

		if !existingSequences[seq] { // Verifica si ya existe
			existingSequences[seq] = true
			return seq
		}
	}
}

func getExistingSequences(db *sql.DB, tornament bool) map[string]bool {
	sequences := make(map[string]bool)

	rows, err := func() (*sql.Rows, error) {
		if tornament {
			return db.Query("SELECT code_id FROM torneos")
		}
		return db.Query("SELECT friend_id FROM user_basic_info")
	}()

	if err != nil {
		log.Fatal("Error al obtener secuencias:", err)
	}
	defer rows.Close()

	for rows.Next() {
		var seq string
		if err := rows.Scan(&seq); err != nil {
			log.Fatal("Error escaneando secuencia:", err)
		}
		sequences[seq] = true
	}
	return sequences
}


