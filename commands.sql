-- Effacer le contenu des tables
-- DELETE FROM posts;
DELETE FROM users;

-- Effacer les tables
-- DROP TABLE users;

CREATE TABLE IF NOT EXISTS users (username TEXT NOT NULL UNIQUE, email TEXT NOT NULL UNIQUE, password TEXT NOT NULL UNIQUE)

