FROM golang:1.23.6-alpine

WORKDIR /app

# Instalar dependencias del sistema
RUN apk add --no-cache git

# Copiar archivos de dependencias
COPY go.mod go.sum ./

# Descargar dependencias
RUN go mod download

# Copiar el código fuente
COPY . .

# Compilar la aplicación
RUN go build -o main ./cmd/api

# Exponer el puerto
EXPOSE 9001

# Comando para ejecutar la aplicación
CMD ["./main"] 