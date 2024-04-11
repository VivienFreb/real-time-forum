-- CREATE TABLE IF NOT EXISTS `Users` (
--     id INTEGER PRIMARY KEY AUTOINCREMENT,
--     pseudo VARCHAR(25) NOT NULL UNIQUE,
--     email  VARCHAR(45) NOT NULL UNIQUE,
--     password TEXT NOT NULL
-- );

-- CREATE TABLE IF NOT EXISTS Users (username TEXT, email TEXT, password TEXT )

-- DElETE FROM `Users`
-- INSERT INTO Posts (title,description,user_id) 
-- VALUES (
--     "Cookies",
--     "J'aime les cookies. C'est bon les cookies. Sauf ceux entièrement au chocolat.",
--     (SELECT id FROM Users WHERE username = "DaBoi27")
--     );

-- INSERT INTO `Disscussions`(speaker, listener, content) VALUES ("Bento-chan","DaBoi27","Dans le pire des cas le McDo m'a donné une fiche de recrutement...")
-- INSERT INTO `Disscussions`(speaker, listener, content) VALUES ("DaBoi27","Bento-chan","Parce qu'ils recrutent des codeurs?")
-- INSERT INTO `Disscussions`(speaker, listener, content) VALUES ("Bento-chan","DaBoi27","Apparement!")
-- INSERT INTO `Disscussions`(speaker, listener, content) VALUES ("DaBoi27","Naruto","Je serais le Roi des Pirates Saiyan Hokage du Dragon de Feu des Jedis de la Communauté de l'Anneau affilié à Poudlard!")


-- CREATE TABLE Test (time CURRENT_TIMESTAMP)
INSERT INTO Test (time) VALUES ("2006-01-02 15:04:05")