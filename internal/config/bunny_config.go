package config

import (
	"fmt"
	"os"

	"git.sr.ht/~jamesponddotco/bunnystorage-go"
)

// InitBunnyStorageClient inicializa y devuelve un cliente de BunnyStorage configurado con credenciales
// desde variables de entorno
func InitBunnyStorageClient() (*bunnystorage.Client, string, error) {
	// Obtener credenciales desde variables de entorno
	readOnlyKey := os.Getenv("BUNNYNET_READ_API_KEY")
	readWriteKey := os.Getenv("BUNNYNET_WRITE_API_KEY")
	storageZone := os.Getenv("BUNNYNET_STORAGE_ZONE")

	// Verificar que las credenciales estén disponibles
	if readOnlyKey == "" || readWriteKey == "" || storageZone == "" {
		return nil, "", fmt.Errorf("credenciales de BunnyStorage no configuradas en variables de entorno")
	}

	endpoint := "falkenstein" // Endpoint por defecto

	// Determinar el endpoint correcto según el valor configurado
	var bunnyEndpoint bunnystorage.Endpoint
	switch endpoint {
	case "falkenstein":
		bunnyEndpoint = bunnystorage.EndpointFalkenstein
	case "stockholm":
		bunnyEndpoint = bunnystorage.EndpointStockholm
	// Nota: Algunos endpoints pueden no estar disponibles en la versión actual de la biblioteca
	// Se usa Falkenstein como valor predeterminado para estos casos
	default:
		bunnyEndpoint = bunnystorage.EndpointFalkenstein
	}

	// Crear la configuración para BunnyStorage
	cfg := &bunnystorage.Config{
		StorageZone: storageZone,
		Key:         readWriteKey,
		ReadOnlyKey: readOnlyKey,
		Endpoint:    bunnyEndpoint,
	}

	// Crear el cliente con la configuración
	client, err := bunnystorage.NewClient(cfg)
	if err != nil {
		return nil, "", fmt.Errorf("error al crear el cliente de BunnyStorage: %v", err)
	}

	return client, storageZone, nil
}
