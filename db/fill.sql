-- database: ./forum.sqlite

-- Users
INSERT INTO `Users`(username , email , password, status) VALUES ("DaBoi27","DaBoi@gmx.fr","DaBoi10@","inactive");
INSERT INTO `Users`(username , email , password, status) VALUES ("Bento-chan","shizuka@gmx.fr","Bento10@","inactive");
-- Posts
INSERT INTO Posts (title,description,user_id) 
VALUES (
    "Cookies",
    "J'aime les cookies. C'est bon les cookies. Sauf ceux entièrement au chocolat.",
    (SELECT id FROM Users WHERE username = "DaBoi27")
    );
-- INSERT INTO Posts(title,imgUrl,description,date,user_id) VALUES ("Lasagne","/asset/image/images.png","Lasagne de mamie","2018-05-20 19:21:00",(SELECT id FROM Users WHERE pseudo = "admin"));
-- INSERT INTO Posts(title,imgUrl,description,date,user_id) VALUES ("Soupe","/asset/image/images.png","Soupe de mamie","2020-05-20 19:21:00",(SELECT id FROM Users WHERE pseudo = "admin2"));