# Manual Técnico - Proyecto Verde

## Índice
1. [Introducción](#introducción)
2. [Arquitectura del Sistema](#arquitectura-del-sistema)
3. [Componentes del Backend](#componentes-del-backend)
4. [Frontend Flutter Web](#frontend-flutter-web)
5. [Base de Datos](#base-de-datos)
6. [Configuración del Entorno](#configuración-del-entorno)
7. [Despliegue](#despliegue)
8. [Flujos de Trabajo Principales](#flujos-de-trabajo-principales)
9. [Procedimientos de Desarrollo](#procedimientos-de-desarrollo)
10. [Integración con Servicios Externos](#integración-con-servicios-externos)
11. [Seguridad](#seguridad)
12. [Mantenimiento](#mantenimiento)
13. [Solución de Problemas Comunes](#solución-de-problemas-comunes)

## Introducción
Este manual técnico documenta la arquitectura, componentes y procedimientos para el mantenimiento y despliegue del Proyecto Verde. El documento está dirigido a desarrolladores y administradores de sistemas que necesitan comprender los aspectos técnicos de la aplicación.

## Arquitectura del Sistema
### Visión General
El sistema está compuesto por un backend desarrollado en Go, una base de datos PostgreSQL, y se despliega utilizando Docker y Traefik como reverse proxy. La aplicación implementa una API RESTful para la comunicación con el frontend.

### Estructura de Directorios
```
proyecto-verde/
├── cmd/
│   └── api/            # Punto de entrada de la aplicación
├── internal/           # Código interno de la aplicación
│   ├── config/         # Configuración de la aplicación
│   ├── handlers/       # Manejadores HTTP
│   ├── middleware/     # Middleware de la aplicación
│   ├── repository/     # Acceso a datos
│   └── routes/         # Definición de rutas
├── pkg/                # Paquetes reutilizables
│   └── database/       # Utilidades de base de datos
├── db/                 # Scripts de base de datos
├── Dockerfile          # Configuración de Docker
└── docker-compose.yml  # Configuración de servicios
```

## Componentes del Backend
### Estructura del Código
El backend está desarrollado en Go y sigue una arquitectura modular. A continuación se detallan los principales componentes:

#### Punto de Entrada (main.go)
El archivo `cmd/api/main.go` inicializa la aplicación, conecta a la base de datos y configura los servicios necesarios. Sus principales funciones son:

- Cargar variables de entorno desde el archivo `.env`
- Establecer la conexión con la base de datos PostgreSQL
- Inicializar el cliente de BunnyStorage para el almacenamiento de archivos
- Configurar los repositorios y handlers
- Configurar las rutas de la API
- Configurar CORS para permitir solicitudes desde diferentes orígenes
- Iniciar el servidor HTTP en el puerto especificado

```go
func main() {
    // Cargar variables de entorno
    // Configurar la base de datos
    // Conectar a la base de datos
    // Inicializar el cliente de BunnyStorage
    // Inicializar repositorios
    // Inicializar handlers
    // Configurar rutas
    // Configurar CORS
    // Iniciar servidor
}
```

#### Sistema de Rutas
Las rutas de la API están definidas en `internal/routes/routes.go` utilizando el router Gorilla Mux. El sistema de rutas organiza todos los endpoints de la aplicación en grupos lógicos:

- Rutas de usuario: gestión de cuentas, perfil y estadísticas
- Rutas de ranking: listado de usuarios clasificados
- Rutas de torneos: creación, gestión e inscripción a torneos
- Rutas de acciones de usuario: registro de actividades
- Rutas de amigos: sistema de gestión de amistades
- Rutas de medallas: sistema de logros y recompensas

Todas las rutas siguen el prefijo `/api/` para las peticiones directas al backend y pueden también ser accedidas mediante el prefijo `/v1/` cuando se utilizan a través del proxy.

#### Controladores (Handlers)
Los controladores manejan las solicitudes HTTP y se encuentran en `internal/handlers/`. Cada módulo funcional tiene su propio handler:

- `UserHandler`: gestión de usuarios y perfiles
- `TorneoHandler`: gestión de torneos y participaciones
- `UserActionsHandler`: registro de acciones de los usuarios
- `UserFriendsHandler`: sistema de amistad entre usuarios
- `MedallasHandler`: gestión del sistema de medallas y logros

#### Repositorios
La capa de acceso a datos está implementada en `internal/repository/postgres/`. Cada entidad principal tiene su propio repositorio:

- `UserRepository`: operaciones CRUD para usuarios
- `TorneoRepository`: operaciones para torneos
- `UserActionsRepository`: registro de acciones
- `UserFriendsRepository`: gestión de relaciones entre usuarios
- `MedallasRepository`: gestión de medallas y asignación a usuarios

### API RESTful
La aplicación expone una API RESTful con los siguientes endpoints principales:

#### Gestión de Usuarios
- `POST /api/users`: Crear nuevo usuario
- `POST /api/auth/login`: Autenticar usuario
- `GET /api/auth/relogin/{id}`: Re-autenticar usuario por ID
- `GET /api/users/{id}`: Obtener información de usuario
- `PUT /api/users/{id}`: Actualizar información de usuario
- `POST/PUT /api/users/{id}/basic-info`: Gestionar información básica
- `GET /api/users/{id}/profile`: Obtener perfil completo
- `PUT /api/users/{id}/profile/edit`: Editar perfil
- `GET/PUT /api/users/{id}/stats`: Gestionar estadísticas

#### Ranking
- `GET /api/ranking`: Obtener ranking general
- `GET /api/ranking/torneo/{torneo_id}`: Obtener ranking de un torneo

#### Torneos
- `POST /api/torneos`: Crear torneo
- `GET /api/torneos`: Listar torneos
- `GET /api/torneos/{id}`: Obtener información de torneo
- `GET /api/torneos/code/{code_id}`: Obtener torneo por código
- `GET /api/torneos/admin/{id}`: Obtener información de administración
- `POST /api/torneos/admin/{id}/terminar`: Finalizar torneo
- `POST /api/torneos/admin/{id}/borrar`: Eliminar torneo
- `POST /api/torneos/inscribir/{code_id}`: Inscribir usuario en torneo
- `DELETE /api/torneos/{torneo_id}/usuario/{user_id}`: Salir de torneo
- `PUT /api/torneos/{id}`: Actualizar torneo
- `GET /api/torneos/{id}/estadisticas`: Obtener estadísticas de torneo
- `GET /api/users/{user_id}/torneos`: Obtener torneos de usuario
- `GET /api/torneos/{torneo_id}/usuario/{user_id}/equipo`: Obtener equipo de usuario en torneo

#### Acciones de Usuario
- `POST /api/users/{user_id}/actions`: Registrar acción
- `GET /api/users/{user_id}/actions`: Obtener acciones de usuario
- `DELETE /api/actions/{id}`: Eliminar acción
- `GET /api/actions`: Obtener todas las acciones

#### Sistema de Amigos
- `GET /api/users/{user_id}/friends`: Listar amigos
- `POST /api/users/{user_id}/friends/add`: Enviar solicitud de amistad
- `PUT /api/users/{user_id}/friends/{friend_id}/accept`: Aceptar solicitud
- `DELETE /api/users/{user_id}/friends/{friend_id}`: Eliminar amigo

#### Sistema de Medallas
- `POST /api/medallas`: Crear medalla
- `GET /api/medallas`: Listar medallas
- `GET /api/users/{user_id}/medallas`: Obtener medallas de usuario
- `POST /api/users/{user_id}/medallas/{medalla_id}`: Asignar medalla
- `GET /api/users/{user_id}/medallas/slogans`: Obtener slogans de medallas
- `GET /api/users/{user_id}/medallas/reset-pending`: Resetear medallas pendientes

## Frontend Flutter Web
### Descripción General
El frontend de la aplicación está desarrollado con Flutter Web, lo que permite tener una aplicación web progresiva (PWA) con una experiencia de usuario fluida y moderna. Flutter Web utiliza tecnologías web estándar (HTML, CSS, JavaScript) pero aprovecha el framework de Flutter para ofrecer una experiencia de desarrollo unificada.

### Estructura del Proyecto Flutter
La estructura típica del proyecto Flutter para el frontend incluye:

```
frontend/
├── lib/
│   ├── main.dart              # Punto de entrada de la aplicación
│   ├── app/                   # Configuración de la aplicación
│   ├── config/                # Configuraciones y constantes
│   ├── models/                # Modelos de datos
│   ├── providers/             # Gestores de estado (si se usa Provider)
│   ├── repositories/          # Acceso a API y datos remotos
│   ├── screens/               # Pantallas de la aplicación
│   ├── services/              # Servicios de la aplicación
│   ├── utils/                 # Utilidades y helpers
│   └── widgets/               # Widgets reutilizables
├── assets/                    # Recursos estáticos (imágenes, fuentes)
├── web/                       # Configuración específica para web
│   ├── index.html             # Página HTML principal
│   ├── manifest.json          # Configuración de PWA
│   └── favicon.ico            # Icono de la aplicación
├── pubspec.yaml               # Dependencias del proyecto
└── test/                      # Pruebas automatizadas
```

### Proceso de Compilación
El frontend se compila usando el siguiente comando específico para generar código optimizado con soporte para WebAssembly:

```bash
flutter build web --release --wasm
```

Este comando:
1. Compila el código Dart a WebAssembly (WASM) para mejor rendimiento
2. Optimiza los recursos para producción
3. Genera los archivos estáticos necesarios para desplegar la aplicación web

Los archivos compilados se generan en el directorio `build/web/` y contienen todo lo necesario para desplegar la aplicación, incluyendo:
- Archivos HTML, CSS y JavaScript
- Assets compilados y optimizados
- Código WASM para mejor rendimiento
- Configuración de Service Worker para funcionalidades de PWA

### Integración con el Backend
El frontend se comunica con el backend a través de la API RESTful usando:

1. **HTTP Client**: Flutter utiliza un cliente HTTP para realizar peticiones a los endpoints del backend.
2. **Modelos de Datos**: Los modelos en Flutter reflejan las estructuras de datos del backend.
3. **Gestión de Estado**: Implementa gestores de estado para manejar la información de la aplicación.

Ejemplo de comunicación típica:

```dart
// Ejemplo simplificado de un servicio de API en Flutter
class ApiService {
  final String baseUrl = 'https://vive.integra-expansion.com/v1/api';
  
  Future<User> loginUser(String email, String password) async {
    final response = await http.post(
      Uri.parse('$baseUrl/auth/login'),
      body: jsonEncode({
        'email': email,
        'password': password,
      }),
      headers: {'Content-Type': 'application/json'},
    );
    
    if (response.statusCode == 200) {
      return User.fromJson(jsonDecode(response.body));
    } else {
      throw Exception('Error de autenticación');
    }
  }
}
```

### Despliegue del Frontend
En el entorno de producción, los archivos compilados de Flutter Web se sirven a través del contenedor NGINX definido en el archivo `docker-compose.yml`. Este servidor web es responsable de entregar el contenido estático a los usuarios y está configurado para funcionar con Traefik como reverse proxy.

## Base de Datos
Se utiliza PostgreSQL como sistema de gestión de base de datos relacional.

### Esquema
El esquema de la base de datos está definido en el archivo `db/schema.sql` y se inicializa automáticamente durante el despliegue con Docker. El sistema utiliza UUID como identificadores primarios con la extensión `uuid-ossp` de PostgreSQL.

### Modelo de Datos
El sistema utiliza un modelo relacional con las siguientes tablas principales:

#### Usuarios y Perfiles
- **user_access**: Almacena la información de autenticación de usuarios.
  - `id`: UUID único del usuario (PK, generado automáticamente)
  - `username`: Nombre de usuario (único)
  - `password`: Contraseña del usuario (debe almacenarse hasheada)

- **user_profile**: Contiene información de personalización del avatar del usuario.
  - `id`: UUID único del perfil (PK)
  - `user_id`: Referencia al usuario (FK, único)
  - `slogan`: Frase o lema del usuario
  - `cabello`: Estilo de cabello seleccionado
  - `vestimenta`: Tipo de vestimenta seleccionada
  - `barba`: Estilo de barba seleccionado
  - `detalle_facial`: Detalles faciales adicionales
  - `detalle_adicional`: Otros detalles de personalización

- **user_basic_info**: Contiene información básica del usuario.
  - `id`: UUID único (PK)
  - `user_id`: Referencia al usuario (FK, único)
  - `numero`: Número telefónico (único)
  - `nombre`: Nombre del usuario
  - `apellido`: Apellido del usuario
  - `friend_id`: Identificador único para amistades (único)

- **user_stats**: Estadísticas acumuladas de los usuarios.
  - `id`: UUID único (PK)
  - `user_id`: Referencia al usuario (FK, único)
  - `puntos`: Puntos totales acumulados
  - `acciones`: Número de acciones realizadas
  - `torneos_participados`: Número de torneos en los que ha participado
  - `torneos_ganados`: Número de torneos ganados
  - `cantidad_amigos`: Número de amigos
  - `es_dueno_torneo`: Indica si el usuario ha creado torneos
  - `pending_medalla`: Número de medallas pendientes de revisar
  - `pending_amigo`: Número de solicitudes de amistad pendientes
  - `torneo_id`: Torneo actual (opcional, FK)

#### Torneos
- **torneos**: Información de los torneos.
  - `id`: UUID único del torneo (PK)
  - `id_creator`: Usuario creador (FK)
  - `nombre`: Nombre del torneo
  - `modalidad`: Modalidad del torneo ('Versus' o 'Individual')
  - `ubicacion_a_latitud`: Latitud de la primera ubicación
  - `ubicacion_a_longitud`: Longitud de la primera ubicación
  - `nombre_ubicacion_a`: Nombre de la primera ubicación
  - `ubicacion_b_latitud`: Latitud de la segunda ubicación (opcional)
  - `ubicacion_b_longitud`: Longitud de la segunda ubicación (opcional)
  - `nombre_ubicacion_b`: Nombre de la segunda ubicación (opcional)
  - `fecha_inicio`: Fecha y hora de inicio
  - `fecha_fin`: Fecha y hora de finalización
  - `ubicacion_aproximada`: Indica si se permite ubicación aproximada
  - `metros_aproximados`: Margen de error en metros para la aproximación
  - `finalizado`: Estado de finalización del torneo
  - `code_id`: Código único para unirse al torneo
  - `ganador_versus`: Equipo ganador en modalidad versus (booleano)
  - `ganador_individual`: Usuario ganador en modalidad individual (UUID)

- **torneo_estadisticas**: Estadísticas de participantes en torneos.
  - `id`: UUID único (PK)
  - `id_jugador`: Referencia al usuario participante (FK)
  - `equipo`: Equipo asignado (booleano)
  - `id_torneo`: Referencia al torneo (FK)
  - `modalidad`: Modalidad de participación
  - `puntos`: Puntos acumulados en el torneo
  - `habilitado`: Estado de participación

#### Acciones
- **user_actions**: Registro de actividades de los usuarios.
  - `id`: UUID único de la acción (PK)
  - `user_id`: Referencia al usuario (FK)
  - `tipo_accion`: Tipo de acción ('ayuda', 'alerta', 'descubrimiento')
  - `foto`: URL de la foto registrada
  - `latitud`: Latitud de la ubicación
  - `longitud`: Longitud de la ubicación
  - `ciudad`: Ciudad donde se realizó la acción
  - `lugar`: Lugar específico de la acción
  - `en_colaboracion`: Indica si se realizó con otros usuarios
  - `colaboradores`: Array de UUIDs de colaboradores
  - `es_para_torneo`: Indica si cuenta para un torneo
  - `id_torneo`: Torneo asociado (opcional, FK)
  - `created_at`: Fecha y hora de creación
  - `deleted_at`: Fecha y hora de eliminación (para borrado lógico)

#### Amistades
- **user_friends**: Sistema de amistad entre usuarios.
  - `id`: UUID único (PK)
  - `user_id`: Referencia al usuario solicitante (FK, parte de PK compuesta)
  - `friend_id`: Referencia al usuario amigo (FK, parte de PK compuesta)
  - `pending_id`: Usuario pendiente de aceptación (FK, opcional)
  - `deleted_at`: Fecha y hora de eliminación (para borrado lógico)
  - `created_at`: Fecha y hora de creación
  - Restricción para evitar auto-amistades

#### Medallas y Logros
- **medallas**: Catálogo de medallas disponibles.
  - `id`: UUID único de la medalla (PK)
  - `nombre`: Nombre de la medalla
  - `descripcion`: Descripción de la medalla
  - `dificultad`: Nivel de dificultad (1-4)
  - Campos de requisitos:
    - `requiere_amistades`: Requiere cierto número de amigos
    - `requiere_puntos`: Requiere cierta cantidad de puntos
    - `requiere_acciones`: Requiere realizar cierto número de acciones
    - `requiere_torneos`: Requiere participar en cierto número de torneos
    - `requiere_victoria_torneos`: Requiere ganar cierto número de torneos
  - `numero_requerido`: Cantidad necesaria para obtener la medalla

- **medallas_ganadas**: Relación entre usuarios y medallas obtenidas.
  - `id`: UUID único (PK)
  - `id_usuario`: Referencia al usuario (FK)
  - `id_medalla`: Referencia a la medalla (FK)
  - `fecha_ganada`: Fecha y hora de obtención

### Relaciones Clave
- Un usuario (`user_access`) tiene información básica (`user_basic_info`), perfil (`user_profile`) y estadísticas (`user_stats`)
- Un usuario puede crear múltiples torneos (`torneos`)
- Un usuario puede participar en múltiples torneos (`torneo_estadisticas`)
- Un usuario puede registrar múltiples acciones (`user_actions`)
- Un usuario puede tener múltiples amigos (`user_friends`)
- Un usuario puede ganar múltiples medallas (`medallas_ganadas`)
- Las acciones pueden estar asociadas a torneos específicos (`user_actions` -> `torneos`)

### Índices y Optimización
Para mejorar el rendimiento de las consultas más frecuentes, el sistema aprovecha:

1. Claves primarias basadas en UUID que se generan automáticamente
2. Restricciones de unicidad en campos críticos como `username`, `numero` y `friend_id`
3. Índices implícitos en claves foráneas para búsquedas rápidas
4. Restricciones CHECK para validar datos (por ejemplo, en tipos de acción)
5. Clave primaria compuesta en la tabla `user_friends` para optimizar consultas de amistad

### Datos Precargados
El esquema incluye la inserción inicial de medallas con diferentes requisitos y niveles de dificultad:

1. Medallas basadas en amistades:
   - "Vengadores, ¡unidos!" (10 amigos)
   - "Familia" (5 amigos)
   - "El poder de la amistad" (20 amigos)

2. Medallas basadas en puntos:
   - "Making my way downtown" (100 puntos)
   - "This is Sparta!" (300 puntos)
   - "Soy inevitable" (1000 puntos)

3. Medallas basadas en acciones:
   - "Orden 66" (66 acciones)
   - "O limpias, o te limpio." (500 acciones)

4. Medallas basadas en torneos:
   - "Que comience el juego" (participar en 1 torneo)
   - "El elegido" (participar en 10 torneos)

5. Medallas basadas en victorias:
   - "No hay nadie que me gane" (ganar 3 torneos)
   - "Super Saiyajin" (ganar 5 torneos seguidos)

### Transacciones y Consistencia
El sistema utiliza transacciones para operaciones críticas como:
- Registro de acciones que otorgan puntos
- Asignación de medallas
- Inscripción en torneos

### Migración y Actualizaciones
Para actualizar el esquema de la base de datos:
1. Crear scripts de migración con cambios incrementales
2. Probar en entorno de desarrollo
3. Aplicar con tiempo de inactividad planificado
4. Tener script de rollback preparado

## Configuración del Entorno
### Variables de Entorno
La aplicación utiliza variables de entorno para la configuración. Estas variables se definen en archivos `.env` y son cargadas mediante la biblioteca `godotenv`. Las principales variables de configuración son:

#### Configuración de Base de Datos
- `DB_HOST`: Host de la base de datos (por defecto: "localhost")
- `DB_PORT`: Puerto de la base de datos (por defecto: "5432")
- `DB_USER`: Usuario de la base de datos (por defecto: "postgres")
- `DB_PASSWORD`: Contraseña de la base de datos
- `DB_NAME`: Nombre de la base de datos (por defecto: "proyecto_verde")

#### Configuración del Servidor
- `PORT`: Puerto en el que se ejecutará el servidor (por defecto: "9001")
- `CDN_URL`: URL de la CDN para servir archivos estáticos

#### Configuración de BunnyStorage
- `BUNNYNET_READ_API_KEY`: Clave de API de BunnyStorage para lectura
- `BUNNYNET_WRITE_API_KEY`: Clave de API de BunnyStorage para escritura
- `BUNNYNET_STORAGE_ZONE`: Zona de almacenamiento de BunnyStorage
- `BUNNYNET_PULL_ZONE`: Zona de pull de BunnyStorage

### Archivos de Configuración
Existen dos archivos de configuración `.env` en el proyecto:
1. El archivo `.env` en la raíz del proyecto, utilizado para desarrollo local
2. El archivo `cmd/api/.env` que contiene la configuración específica para el entorno de producción

**Importante**: Para el despliegue en producción, es necesario asegurarse de que las claves de API de BunnyStorage estén correctamente configuradas y que la zona de almacenamiento esté disponible.

## Despliegue
### Docker
La aplicación se despliega utilizando Docker y Docker Compose. El archivo `Dockerfile` define la imagen base para el backend:

```dockerfile
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
```

### Docker Compose
El archivo `docker-compose.yml` define los servicios necesarios para ejecutar la aplicación:

1. **reverse-proxy**: Servicio de Traefik para enrutar el tráfico y gestionar SSL
2. **backend**: Servicio que ejecuta la aplicación Go
3. **frontend**: Servicio que sirve los archivos estáticos del frontend mediante NGINX, aquí se encuentra el compilado de la aplicación Flutter Web
4. **postgres**: Servicio de base de datos PostgreSQL

Para desplegar la aplicación en producción, ejecutar:
```bash
docker-compose up -d
```

### Proceso de Compilación y Despliegue del Frontend
Para compilar y desplegar el frontend en la aplicación, se deben seguir estos pasos:

1. **Compilación del Frontend Flutter Web**:
   ```bash
   # Navegar al directorio del frontend
   cd frontend

   # Asegurarse de tener las dependencias actualizadas
   flutter pub get

   # Compilar para producción con soporte WASM para mejor rendimiento
   flutter build web --release --wasm
   ```

2. **Colocar los archivos compilados**:
   Los archivos generados se encontrarán en `frontend/build/web/`. Estos archivos deben copiarse al directorio que NGINX sirve en el contenedor.

   ```bash
   # Copiar los archivos compilados al directorio frontend en la raíz del proyecto
   cp -r build/web/* ../../frontend/
   ```

3. **Configuración de NGINX**:
   El servicio frontend utiliza NGINX para servir los archivos estáticos. La configuración básica se encuentra en `nginx.conf`:

   ```nginx
   server {
       listen 80;
       root /usr/share/nginx/html;
       index index.html;
       
       location / {
           try_files $uri $uri/ /index.html;
       }
       
       # Configuración de caché para recursos estáticos
       location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg|wasm)$ {
           expires 1y;
           add_header Cache-Control "public, max-age=31536000";
       }
   }
   ```

4. **Actualización del Frontend en Producción**:
   Para actualizar solo el frontend en un sistema en producción:

   ```bash
   # Copiar los nuevos archivos compilados al servidor
   scp -r frontend/* usuario@servidor:/ruta/al/proyecto/frontend/

   # O, si ya estamos en el servidor, reemplazar los archivos y reiniciar solo el contenedor frontend
   docker-compose restart frontend

   # O, usar el sistema Git para actualizar el frontend
   git add .
   git commit -m "Actualización del frontend"
   git push
   ```

### Traefik
Se utiliza Traefik como reverse proxy para gestionar el tráfico y configurar SSL. La configuración se encuentra en el archivo `docker-compose.yml` y define:

- Certificados SSL automáticos con Let's Encrypt
- Redirección de HTTP a HTTPS
- Enrutamiento de tráfico entre backend y frontend
- Enrutamiento basado en prefijos de URL (/v1 para backend)

## Flujos de Trabajo Principales
Esta sección describe los principales flujos de trabajo de la aplicación para comprender mejor cómo interactúan los diferentes componentes del sistema.

### Registro y Autenticación de Usuarios
1. **Registro de Usuario**:
   - El usuario completa el formulario de registro en el frontend
   - La solicitud es enviada al endpoint `POST /api/users`
   - El backend valida los datos y crea un nuevo registro en la tabla `users`
   - Se inicializan las tablas relacionadas (`user_profiles`, `user_stats`)
   - Se devuelve un token de autenticación y la información del usuario

2. **Inicio de Sesión**:
   - El usuario introduce credenciales en el formulario de login
   - El frontend envía la solicitud a `POST /api/auth/login`
   - El backend verifica las credenciales contra la base de datos
   - Si son correctas, se genera un token y se devuelve al cliente
   - El token se almacena en el almacenamiento local del navegador

3. **Re-autenticación**:
   - Al cargar la aplicación, se verifica si existe un token guardado
   - Si existe, se realiza una petición a `GET /api/auth/relogin/{id}`
   - El backend valida el token y devuelve los datos actualizados del usuario

### Creación y Gestión de Torneos
1. **Creación de Torneo**:
   - Un usuario administrador crea un torneo mediante el formulario correspondiente
   - El frontend envía los datos a `POST /api/torneos`
   - El backend genera un código único para el torneo y lo almacena en la tabla `torneos`
   - Se devuelve la información del torneo creado, incluyendo el código de acceso

2. **Unirse a un Torneo**:
   - El usuario introduce el código del torneo en la aplicación
   - Se realiza una petición a `GET /api/torneos/code/{code_id}` para verificar el torneo
   - Si el torneo existe y está activo, se muestra la información y opción de unirse
   - Al confirmar, se envía una solicitud a `POST /api/torneos/inscribir/{code_id}`
   - El backend registra al usuario en el torneo y actualiza la tabla `torneo_participants`

3. **Gestión de Torneos**:
   - Los administradores pueden ver sus torneos mediante `GET /api/torneos/admin/{id}`
   - Pueden finalizar un torneo con `POST /api/torneos/admin/{id}/terminar`
   - Pueden eliminar un torneo con `POST /api/torneos/admin/{id}/borrar`
   - Los usuarios pueden ver sus torneos activos con `GET /api/users/{user_id}/torneos`

### Registro de Acciones
1. **Registrar una Acción**:
   - El usuario realiza una actividad y la registra a través de la interfaz
   - Opcionalmente, puede subir una imagen/video como evidencia
   - Si hay multimedia, primero se sube a BunnyStorage
   - La solicitud se envía a `POST /api/users/{user_id}/actions`
   - El backend registra la acción, asigna puntos y actualiza las estadísticas
   - Se verifica si la acción cumple requisitos para obtener medallas

2. **Visualización de Acciones**:
   - Las acciones de un usuario se obtienen mediante `GET /api/users/{user_id}/actions`
   - Se muestran en el perfil del usuario y en el feed de actividad
   - Las imágenes/videos se sirven a través de la CDN de BunnyStorage

### Sistema de Medallas
1. **Asignación Automática de Medallas**:
   - Al registrar acciones, el backend verifica si se cumplen requisitos para medallas
   - Si se cumplen, se asigna automáticamente la medalla al usuario
   - Se marca como pendiente hasta que el usuario la visualice

2. **Asignación Manual**:
   - Un administrador puede asignar medallas manualmente
   - Utiliza el endpoint `POST /api/users/{user_id}/medallas/{medalla_id}`
   - El sistema registra la asignación y notifica al usuario

3. **Visualización de Medallas**:
   - El usuario puede ver sus medallas mediante `GET /api/users/{user_id}/medallas`
   - Las medallas pendientes se muestran destacadas hasta que se visualizan
   - El usuario puede ver los slogans de sus medallas con `GET /api/users/{user_id}/medallas/slogans`

### Sistema de Amigos
1. **Envío de Solicitud de Amistad**:
   - El usuario busca a otro usuario y envía solicitud
   - Se realiza una petición a `POST /api/users/{user_id}/friends/add`
   - El backend registra la solicitud con estado "pendiente"

2. **Aceptación de Solicitud**:
   - El usuario recibe notificación de solicitud pendiente
   - Al aceptar, se llama a `PUT /api/users/{user_id}/friends/{friend_id}/accept`
   - El backend actualiza el estado a "aceptada" y crea la relación bidireccional

3. **Visualización de Amigos**:
   - La lista de amigos se obtiene mediante `GET /api/users/{user_id}/friends`
   - Se muestra en el perfil y en la sección de amigos
   - Se puede filtrar por estado (pendiente/aceptada)

## Procedimientos de Desarrollo
Esta sección describe los procedimientos recomendados para el desarrollo y mantenimiento del código del proyecto.

### Configuración del Entorno de Desarrollo
#### Backend (Go)
1. **Requisitos previos**:
   - Go 1.20 o superior
   - PostgreSQL 15
   - Git

2. **Configuración inicial**:
   ```bash
   # Clonar el repositorio
   git clone https://github.com/tu-organizacion/proyecto-verde.git
   cd proyecto-verde
   
   # Instalar dependencias
   go mod download
   
   # Configurar archivo .env para desarrollo local
   nano .env
   # Editar .env con los valores apropiados
   
   # Inicializar la base de datos
   psql -U postgres -c "CREATE DATABASE proyecto_verde;"
   psql -U postgres -d proyecto_verde -f db/schema.sql
   
   # Ejecutar la aplicación en modo desarrollo
   go run cmd/api/main.go
   ```

#### Frontend (Flutter)
1. **Requisitos previos**:
   - Flutter SDK (versión estable más reciente)
   - Dart SDK
   - Editor (VS Code con extensiones Flutter/Dart recomendado)

2. **Configuración inicial**:
   ```bash
   # Navegar al directorio del frontend
   cd frontend
   
   # Instalar dependencias
   flutter pub get
   
   # Configurar variables de entorno
   cp .env.example .env
   # Editar .env con URL del backend local
   
   # Ejecutar en modo desarrollo
   flutter run -d chrome
   ```

### Gestión de Versiones
Se utiliza versionado semántico (SemVer):
- **Mayor (X.0.0)**: Cambios incompatibles con versiones anteriores
- **Menor (0.X.0)**: Nuevas funcionalidades compatibles con versiones anteriores
- **Parche (0.0.X)**: Correcciones de errores compatibles con versiones anteriores

## Integración con Servicios Externos
### BunnyStorage
La aplicación integra BunnyStorage para el almacenamiento de archivos. La configuración se realiza mediante variables de entorno y se inicializa en `config.InitBunnyStorageClient()`.

BunnyStorage se utiliza principalmente para almacenar:
- Imágenes de perfil de usuario
- Archivos relacionados con torneos
- Otros recursos multimedia

## Seguridad
### Autenticación y Autorización
El sistema implementa un mecanismo básico de autenticación mediante credenciales de usuario. Las contraseñas se almacenan de forma segura en la base de datos (se recomienda verificar que estén hasheadas).

### CORS (Cross-Origin Resource Sharing)
La aplicación utiliza el paquete `rs/cors` para configurar CORS y permitir solicitudes desde diferentes orígenes. La configuración actual es:

```go
corsOptions := cors.New(cors.Options{
    AllowedOrigins:   []string{"*"},
    AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
    AllowCredentials: true,
})
```

**Importante para producción**: La configuración actual de CORS permite solicitudes desde cualquier origen (`*`). Para un entorno de producción, se recomienda restringir los orígenes permitidos a dominios específicos.

### HTTPS
En producción, todo el tráfico HTTP es redirigido automáticamente a HTTPS mediante la configuración de Traefik:

```yaml
- "--entrypoints.web.http.redirections.entrypoint.to=websecure"
- "--entrypoints.web.http.redirections.entrypoint.scheme=https"
```

### Certificados SSL
Los certificados SSL son gestionados automáticamente por Traefik utilizando Let's Encrypt:

```yaml
- "--certificatesresolvers.myresolver.acme.email=avn2000inc@gmail.com"
- "--certificatesresolvers.myresolver.acme.storage=/letsencrypt/acme.json"
- "--certificatesresolvers.myresolver.acme.tlschallenge=true"
```

### Recomendaciones de Seguridad Adicionales
Para mejorar la seguridad del sistema en producción, se recomienda:

1. Implementar rate limiting para prevenir ataques de fuerza bruta
2. Añadir middleware de validación de entrada para prevenir ataques de inyección
3. Configurar políticas de Content Security Policy (CSP)
4. Realizar auditorías de seguridad periódicas
5. Mantener todas las dependencias actualizadas
6. Revisar y rotar regularmente las claves de API de BunnyStorage

## Mantenimiento
### Monitorización
Para monitorizar el estado del sistema, se recomienda:

1. Implementar un sistema de monitorización como Prometheus con Grafana
2. Configurar alertas para CPU, memoria, espacio en disco y errores HTTP
3. Monitorizar los logs del sistema utilizando Loki o un sistema similar

### Logs
La aplicación genera logs utilizando el paquete `log` estándar de Go. Se recomienda:

1. Configurar un sistema de centralización de logs como ELK Stack o Graylog
2. Establecer retención de logs acorde a las políticas de la organización
3. Configurar alertas basadas en patrones de logs para detectar problemas

### Respaldos
Para la base de datos PostgreSQL, se recomienda:

1. Configurar respaldos diarios automáticos
2. Almacenar los respaldos en una ubicación externa (preferiblemente cifrados)
3. Verificar periódicamente la restauración de respaldos

```bash
# Ejemplo de script de respaldo para PostgreSQL
pg_dump -U postgres -d proyecto_verde > backup_$(date +%Y%m%d).sql
```

### Actualizaciones
Para actualizar el sistema a una nueva versión:

1. Realizar pruebas en un entorno de staging
2. Hacer respaldo de la base de datos
3. Actualizar el código fuente
4. Reconstruir y reiniciar los contenedores

```bash
# Actualizar y reiniciar servicios
git pull
docker-compose build
docker-compose down
docker-compose up -d
```

### Escalabilidad
Para escalar el sistema ante mayor carga:

1. Escalar horizontalmente el servicio backend
2. Configurar balanceo de carga
3. Considerar la implementación de cachés distribuidas
4. Optimizar consultas de base de datos

## Solución de Problemas Comunes
### Problemas de Conexión a la Base de Datos
Si la aplicación no puede conectarse a la base de datos:

1. Verificar que el servicio de PostgreSQL esté en ejecución
   ```bash
   docker-compose ps postgres
   ```

2. Verificar la configuración de conexión en el archivo `.env`
   ```bash
   cat .env | grep DB_
   ```

3. Intentar conectarse manualmente a la base de datos
   ```bash
   docker-compose exec postgres psql -U postgres -d proyecto_verde
   ```

### Problemas con Traefik
Si hay problemas con el proxy:

1. Verificar los logs de Traefik
   ```bash
   docker-compose logs reverse-proxy
   ```

2. Comprobar que los puertos 80 y 443 estén disponibles
   ```bash
   netstat -tulpn | grep -E '80|443'
   ```

3. Verificar la configuración de los routers en el dashboard de Traefik
   (accesible en http://your-server-ip:8080 si está configurado)

### Errores en el Servicio Backend
Si el servicio backend presenta errores:

1. Verificar los logs del servicio
   ```bash
   docker-compose logs backend
   ```

2. Comprobar la conectividad con BunnyStorage
   ```bash
   curl -I https://storage.bunnycdn.com/
   ```

3. Reiniciar el servicio
   ```bash
   docker-compose restart backend
   ```

### Problemas con los Certificados SSL
Si hay problemas con los certificados SSL:

1. Verificar que el dominio apunta correctamente al servidor
   ```bash
   dig vive.integra-expansion.com
   ```

2. Revisar los logs de Traefik para errores de Let's Encrypt
   ```bash
   docker-compose logs reverse-proxy | grep "acme"
   ```

3. Renovar manualmente los certificados
   ```bash
   docker-compose exec reverse-proxy traefik cert renew
   ```

### Problemas de Permisos en BunnyStorage
Si hay problemas al subir archivos a BunnyStorage:

1. Verificar las claves de API en el archivo `.env`
2. Comprobar que la zona de almacenamiento existe y está correctamente configurada
3. Verificar que las claves tienen permisos de lectura y escritura
4. Revisar los logs de la aplicación para errores específicos relacionados con BunnyStorage

## Configuración para Producción
Para preparar el sistema para producción, se deben realizar los siguientes cambios:

### 1. Configuración de Variables de Entorno
Crear un archivo `.env` para producción con valores seguros:

```
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=<contraseña_segura>
DB_NAME=proyecto_verde
PORT=9001
CDN_URL=https://cdn.integra-expansion.com

# Configuración de BunnyStorage
BUNNYNET_READ_API_KEY=<clave_api_lectura>
BUNNYNET_WRITE_API_KEY=<clave_api_escritura>
BUNNYNET_STORAGE_ZONE=vive
BUNNYNET_PULL_ZONE=cdn.integra-expansion.com
```

### 2. Configuración de CORS
Modificar la configuración de CORS en `cmd/api/main.go`:

```go
corsOptions := cors.New(cors.Options{
    AllowedOrigins:   []string{"https://vive.integra-expansion.com"},
    AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
    AllowCredentials: true,
})
```

### 3. Configuración de Docker Compose
Asegurarse de que el archivo `docker-compose.yml` esté configurado correctamente para producción, especialmente:

- Verificar los volúmenes para persistencia de datos
- Configurar reinicio automático de servicios (`restart: always`)
- Asegurar que los directorios de montaje existan y tengan permisos adecuados

### 4. Configuración de Traefik
Verificar la configuración de Traefik en `docker-compose.yml`:

- Configurar correctamente el correo electrónico para Let's Encrypt
- Verificar los dominios configurados
- Asegurar que las redes estén correctamente definidas

### 5. Seguridad de BunnyStorage
Verificar la configuración de seguridad en BunnyStorage:

- Limitar el acceso a la zona de almacenamiento
- Configurar CORS en el lado de BunnyStorage

## Consideraciones Finales
Este manual técnico debe ser actualizado regularmente conforme el sistema evoluciona. Se recomienda revisar y actualizar la documentación al menos después de cada despliegue importante o cambio significativo en la arquitectura.