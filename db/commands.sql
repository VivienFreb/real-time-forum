-- CREATE TABLE IF NOT EXISTS `Users` (
--     id INTEGER PRIMARY KEY AUTOINCREMENT,
--     pseudo VARCHAR(25) NOT NULL UNIQUE,
--     email  VARCHAR(45) NOT NULL UNIQUE,
--     password TEXT NOT NULL
-- );

-- CREATE TABLE IF NOT EXISTS Users (username TEXT, email TEXT, password TEXT )

-- DElETE FROM `Users`
INSERT INTO Posts (title,description,user_id) 
VALUES (
    "Cookies",
    "J'aime les cookies. C'est bon les cookies. Sauf ceux entièrement au chocolat.",
    (SELECT id FROM Users WHERE username = "DaBoi27")
    );