CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Información de inicio de sesión (tabla base)
CREATE TABLE user_access (
  id UUID NOT NULL DEFAULT uuid_generate_v4(),
  username TEXT NOT NULL UNIQUE,
  password TEXT NOT NULL,
  CONSTRAINT pk_user_access PRIMARY KEY (id)
);

-- Torneos (necesita user_access)
CREATE TABLE torneos (
  id UUID NOT NULL DEFAULT uuid_generate_v4(),
  id_creator UUID NOT NULL,
  nombre VARCHAR(255) NOT NULL,
  modalidad VARCHAR(50) CHECK (modalidad IN ('Versus', 'Individual')) NOT NULL,
  ubicacion_a_latitud DOUBLE PRECISION NOT NULL,
  ubicacion_a_longitud DOUBLE PRECISION NOT NULL,
  nombre_ubicacion_a VARCHAR(255) NOT NULL,
  ubicacion_b_latitud DOUBLE PRECISION,
  ubicacion_b_longitud DOUBLE PRECISION,
  nombre_ubicacion_b VARCHAR(255),
  fecha_inicio TIMESTAMP NOT NULL,
  fecha_fin TIMESTAMP NOT NULL,
  ubicacion_aproximada BOOLEAN NOT NULL DEFAULT FALSE,
  kilometros_aproximados INT,
  finalizado BOOLEAN NOT NULL DEFAULT FALSE,
  code_id TEXT NOT NULL UNIQUE,
  ganador_versus BOOLEAN,
  ganador_individual UUID,
  CONSTRAINT pk_torneos PRIMARY KEY (id),
  CONSTRAINT fk_user_creator_user FOREIGN KEY (id_creator) REFERENCES user_access(id) ON DELETE CASCADE
);

-- Estadísticas de usuario (necesita user_access y torneos)
CREATE TABLE user_stats (
  id UUID NOT NULL DEFAULT uuid_generate_v4(),
  user_id UUID NOT NULL UNIQUE,
  puntos INT NOT NULL DEFAULT 0,
  acciones INT NOT NULL DEFAULT 0,
  torneos_participados INT NOT NULL DEFAULT 0,
  torneos_ganados INT NOT NULL DEFAULT 0,
  cantidad_amigos INT NOT NULL DEFAULT 0,
  es_dueno_torneo BOOLEAN NOT NULL DEFAULT FALSE,
  pending_medalla INT NOT NULL DEFAULT 0,
  pending_amigo INT NOT NULL DEFAULT 0,
  torneo_id UUID,
  CONSTRAINT pk_user_stats PRIMARY KEY (id),
  CONSTRAINT fk_user_stats_user FOREIGN KEY (user_id) REFERENCES user_access(id) ON DELETE CASCADE,
  CONSTRAINT fk_user_stats_torneo FOREIGN KEY (torneo_id) REFERENCES torneos(id) ON DELETE CASCADE
);

-- Información de personalización de personaje (necesita user_access)
CREATE TABLE user_profile (
  id UUID NOT NULL DEFAULT uuid_generate_v4(),
  user_id UUID NOT NULL UNIQUE,
  slogan TEXT,
  cabello TEXT,
  vestimenta TEXT,
  barba TEXT,
  detalle_facial TEXT,
  detalle_adicional TEXT,
  CONSTRAINT pk_user_profile PRIMARY KEY (id),
  CONSTRAINT fk_user_profile_user FOREIGN KEY (user_id) REFERENCES user_access(id) ON DELETE CASCADE
);

-- Información básica del usuario (necesita user_access)
CREATE TABLE user_basic_info (
  id UUID NOT NULL DEFAULT uuid_generate_v4(),
  user_id UUID NOT NULL UNIQUE,
  numero TEXT NOT NULL UNIQUE,
  nombre TEXT NOT NULL,
  apellido TEXT NOT NULL,
  friend_id TEXT NOT NULL UNIQUE,
  CONSTRAINT pk_user_basic_info PRIMARY KEY (id),
  CONSTRAINT fk_user_basic_info_user FOREIGN KEY (user_id) REFERENCES user_access(id) ON DELETE CASCADE
);

-- Amigos (necesita user_access)
CREATE TABLE user_friends (
  id UUID NOT NULL DEFAULT uuid_generate_v4(),
  user_id UUID NOT NULL,
  friend_id UUID NOT NULL,
  pending_id UUID,
  deleted_at TIMESTAMP,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT pk_user_friends PRIMARY KEY (user_id, friend_id),
  CONSTRAINT fk_user_friends_user FOREIGN KEY (user_id) REFERENCES user_access(id) ON DELETE CASCADE,
  CONSTRAINT fk_user_friends_friend FOREIGN KEY (friend_id) REFERENCES user_access(id) ON DELETE CASCADE,
  CONSTRAINT fk_user_friends_pending FOREIGN KEY (pending_id) REFERENCES user_access(id) ON DELETE CASCADE,
  CONSTRAINT check_no_self_friendship CHECK (user_id <> friend_id)
);

-- Estadísticas de torneo (necesita user_access y torneos)
CREATE TABLE torneo_estadisticas (
  id UUID NOT NULL DEFAULT uuid_generate_v4(),
  id_jugador UUID NOT NULL,
  equipo BOOLEAN NOT NULL,
  id_torneo UUID NOT NULL,
  modalidad VARCHAR(50) NOT NULL,
  puntos INT NOT NULL DEFAULT 0,
  habilitado BOOLEAN NOT NULL DEFAULT true,
  CONSTRAINT pk_torneo_participantes PRIMARY KEY (id),
  CONSTRAINT fk_torneo_participantes_jugador FOREIGN KEY (id_jugador) REFERENCES user_access(id) ON DELETE CASCADE,
  CONSTRAINT fk_torneo_participantes_torneo FOREIGN KEY (id_torneo) REFERENCES torneos(id) ON DELETE CASCADE
);

-- Acciones de usuario (necesita user_access y torneos)
CREATE TABLE user_actions (
  id UUID NOT NULL DEFAULT uuid_generate_v4(),
  user_id UUID NOT NULL,
  tipo_accion VARCHAR(50) CHECK (tipo_accion IN ('ayuda', 'alerta', 'descubrimiento')) NOT NULL,
  foto VARCHAR(255) NOT NULL,
  latitud DOUBLE PRECISION NOT NULL,
  longitud DOUBLE PRECISION NOT NULL,
  ciudad VARCHAR(100) NOT NULL,
  lugar VARCHAR(150) NOT NULL,
  en_colaboracion BOOLEAN NOT NULL DEFAULT FALSE,
  colaboradores UUID[],
  es_para_torneo BOOLEAN NOT NULL DEFAULT FALSE,
  id_torneo UUID,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP,
  CONSTRAINT pk_user_actions PRIMARY KEY (id),
  CONSTRAINT fk_user_actions_user FOREIGN KEY (user_id) REFERENCES user_access(id) ON DELETE CASCADE,
  CONSTRAINT fk_user_actions_torneo FOREIGN KEY (id_torneo) REFERENCES torneos(id) ON DELETE CASCADE
);

-- Medallas (tabla independiente)
CREATE TABLE medallas (
  id UUID NOT NULL DEFAULT uuid_generate_v4(),
  nombre VARCHAR(255) NOT NULL,
  descripcion TEXT,
  dificultad INT NOT NULL,
  requiere_amistades BOOLEAN NOT NULL DEFAULT FALSE,
  requiere_puntos BOOLEAN NOT NULL DEFAULT FALSE,
  requiere_acciones BOOLEAN NOT NULL DEFAULT FALSE,
  requiere_torneos BOOLEAN NOT NULL DEFAULT FALSE,
  requiere_victoria_torneos BOOLEAN NOT NULL DEFAULT FALSE,
  numero_requerido INT,
  CONSTRAINT pk_medallas PRIMARY KEY (id)
);

-- Medallas ganadas (necesita user_access y medallas)
CREATE TABLE medallas_ganadas (
  id UUID NOT NULL DEFAULT uuid_generate_v4(),
  id_usuario UUID NOT NULL,
  id_medalla UUID NOT NULL,
  fecha_ganada TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT pk_medallas_ganadas PRIMARY KEY (id),
  CONSTRAINT fk_medallas_ganadas_usuario FOREIGN KEY (id_usuario) REFERENCES user_access(id) ON DELETE CASCADE,
  CONSTRAINT fk_medallas_ganadas_medalla FOREIGN KEY (id_medalla) REFERENCES medallas(id) ON DELETE CASCADE
);

-- Insertar medallas iniciales
INSERT INTO medallas (nombre, descripcion, dificultad, requiere_amistades, numero_requerido) VALUES
('Vengadores, ¡unidos!', 'Consigue más de 10 amigos.', 2, TRUE, 10);

INSERT INTO medallas (nombre, descripcion, dificultad, requiere_puntos, numero_requerido) VALUES
('This is Sparta!', 'Alcanza 300 puntos en total.', 2, TRUE, 300);

INSERT INTO medallas (nombre, descripcion, dificultad, requiere_puntos, numero_requerido) VALUES
('Making my way downtown', 'Alcanza 100 puntos en total.', 1, TRUE, 100);

INSERT INTO medallas (nombre, descripcion, dificultad, requiere_acciones, numero_requerido) VALUES
('Orden 66', 'Realiza 66 acciones en el juego.', 3, TRUE, 66);

INSERT INTO medallas (nombre, descripcion, dificultad, requiere_torneos, numero_requerido) VALUES
('Que comience el juego', 'Participa en tu primer torneo.', 1, TRUE, 1);

INSERT INTO medallas (nombre, descripcion, dificultad, requiere_victoria_torneos, numero_requerido) VALUES
('No hay nadie que me gane', 'Gana 3 torneos.', 3, TRUE, 3);

INSERT INTO medallas (nombre, descripcion, dificultad, requiere_puntos, numero_requerido) VALUES
('Soy inevitable', 'Alcanza 1000 puntos en total.', 3, TRUE, 1000);

INSERT INTO medallas (nombre, descripcion, dificultad, requiere_amistades, numero_requerido) VALUES
('Familia', 'Añade a 5 amigos. Recuerda, la familia es lo más importante.', 1, TRUE, 5);

INSERT INTO medallas (nombre, descripcion, dificultad, requiere_acciones, numero_requerido) VALUES
('O limpias, o te limpio.', 'Realiza 500 acciones en el juego.', 4, TRUE, 500);

INSERT INTO medallas (nombre, descripcion, dificultad, requiere_torneos, numero_requerido) VALUES
('El elegido', 'Participa en 10 torneos.', 4, TRUE, 10);

INSERT INTO medallas (nombre, descripcion, dificultad, requiere_victoria_torneos, numero_requerido) VALUES
('Super Saiyajin', 'Gana 5 torneos seguidos.', 4, TRUE, 5);

INSERT INTO medallas (nombre, descripcion, dificultad, requiere_amistades, numero_requerido) VALUES
('El poder de la amistad', 'Consigue 20 amigos.', 4, TRUE, 20);