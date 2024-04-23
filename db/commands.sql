-- Pour insérer des choses directement dans la database.
-- INSERT INTO `Comments` (parent, content, user_name) VALUES ("Erla Situation","Bof, un de perdu dix de retrouver!","DaBoi27")

-- DELETE FROM `Disscussions`

-- DROP TABLE `Users`

CREATE TABLE IF NOT EXISTS Users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL,
    first_name TEXT,
    last_name TEXT,
    gender TEXT,
    age INTEGER,
    email TEXT,
    password TEXT,
    status TEXT
);

-- UPDATE `Users` SET username = 'Bento-chan' WHERE id = 2

DELETE FROM Comments

-- DaBoi27 DaBoi10@ DaBoi@gmx.fr
-- Bento-chan Bento10@ shizuka@gmx.fr
-- Tok/Italics/Cyrus/Acheld/Tankista/Aromage