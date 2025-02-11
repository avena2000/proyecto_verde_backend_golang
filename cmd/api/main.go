package main

import (
	"backend_proyecto_verde/internal/handlers"
	"backend_proyecto_verde/internal/repository/postgres"
	"backend_proyecto_verde/internal/routes"
	"backend_proyecto_verde/pkg/database"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Cargar variables de entorno
	if err := godotenv.Load(); err != nil {
		log.Printf("No se pudo cargar el archivo .env: %v", err)
	}

	// Configuración de la base de datos
	dbConfig := &database.Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", ""),
		DBName:   getEnv("DB_NAME", "proyecto_verde"),
	}

	// Conectar a la base de datos
	db, err := database.NewConnection(dbConfig)
	if err != nil {
		log.Fatalf("Error al conectar a la base de datos: %v", err)
	}
	defer db.Close()

	// Inicializar repositorios
	userRepo := postgres.NewUserRepository(db)
	torneoRepo := postgres.NewTorneoRepository(db)
	userActionsRepo := postgres.NewUserActionsRepository(db)
	userFriendsRepo := postgres.NewUserFriendsRepository(db)
	medallasRepo := postgres.NewMedallasRepository(db, userRepo)

	// Inicializar handlers
	userHandler := handlers.NewUserHandler(userRepo)
	torneoHandler := handlers.NewTorneoHandler(torneoRepo)
	userActionsHandler := handlers.NewUserActionsHandler(userActionsRepo, medallasRepo)
	userFriendsHandler := handlers.NewUserFriendsHandler(userFriendsRepo)
	medallasHandler := handlers.NewMedallasHandler(medallasRepo)

	// Configurar rutas
	router := routes.SetupRoutes(
		userHandler,
		torneoHandler,
		userActionsHandler,
		userFriendsHandler,
		medallasHandler,
	)

	// Iniciar servidor
	port := getEnv("PORT", "8080")
	log.Printf("Servidor iniciado en el puerto %s", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
} 