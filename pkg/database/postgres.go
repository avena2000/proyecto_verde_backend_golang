package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

func NewConnection(config *Config) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.DBName,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

// WithTransaction ejecuta una función dentro de una transacción
// Si la función devuelve un error, la transacción se revierte
// Si la función no devuelve error, la transacción se confirma
func WithTransaction(db *sql.DB, fn func(*sql.Tx) error) error {
	start := time.Now()
	fmt.Println("=== INICIANDO TRANSACCIÓN ===")

	tx, err := db.Begin()
	if err != nil {
		fmt.Printf("ERROR AL INICIAR TRANSACCIÓN: %v\n", err)
		return err
	}

	// Creamos un proxy para la transacción
	txProxy := &txLogger{tx: tx}

	defer func() {
		if p := recover(); p != nil {
			fmt.Println("=== PANIC DETECTADO, HACIENDO ROLLBACK ===")
			tx.Rollback()
			panic(p) // volver a lanzar el pánico después de hacer rollback
		}
	}()
	// Ejecutamos la función con nuestro proxy
	if err := fn(txProxy.tx); err != nil {
		elapsed := time.Since(start)
		fmt.Printf("=== ERROR DESPUÉS DE %v, HACIENDO ROLLBACK: %v ===\n", elapsed, err)
		tx.Rollback()
		return err
	}

	elapsed := time.Since(start)
	fmt.Printf("=== TRANSACCIÓN EXITOSA DESPUÉS DE %v, HACIENDO COMMIT ===\n", elapsed)
	return tx.Commit()
}

// txLogger es un proxy para sql.Tx que registra las operaciones
type txLogger struct {
	tx *sql.Tx
}
