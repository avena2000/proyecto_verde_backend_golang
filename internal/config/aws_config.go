package config

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

// InitAWSSession inicializa y devuelve una sesión de AWS configurada con credenciales
// desde variables de entorno
func InitAWSSession() (*session.Session, error) {
	// Obtener credenciales desde variables de entorno
	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	region := os.Getenv("AWS_REGION")

	// Verificar que las credenciales estén disponibles
	if accessKey == "" || secretKey == "" {
		return nil, fmt.Errorf("credenciales de AWS no configuradas en variables de entorno")
	}

	// Usar región por defecto si no está configurada
	if region == "" {
		region = "us-east-1" // Región por defecto
	}

	// Crear la sesión de AWS con las credenciales
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
	})

	if err != nil {
		return nil, fmt.Errorf("error al crear la sesión de AWS: %v", err)
	}

	return sess, nil
}
