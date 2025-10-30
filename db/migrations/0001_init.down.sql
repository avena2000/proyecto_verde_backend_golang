-- Borrar datos dependientes y tablas en orden inverso
DROP TABLE IF EXISTS medallas_ganadas;
DROP TABLE IF EXISTS medallas;
DROP TABLE IF EXISTS user_actions;
DROP TABLE IF EXISTS torneo_estadisticas;
DROP TABLE IF EXISTS user_friends;
DROP TABLE IF EXISTS user_basic_info;
DROP TABLE IF EXISTS user_profile;
DROP TABLE IF EXISTS user_stats;
DROP TABLE IF EXISTS torneos;
DROP TABLE IF EXISTS user_access;

-- Quitar extensión (opcional)
DROP EXTENSION IF EXISTS "uuid-ossp";


