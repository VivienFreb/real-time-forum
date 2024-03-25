-- database: ./forum.sqlite

-- Users
INSERT INTO `Users`(pseudo , email , password,access) VALUES ("admin","admin@test.fr","$2a$10$zah2W8Y/jutBoUS/H3Zo.OABhYsC3qFrNesAifp6UK0KgvVj0pSwu",0);
-- mdp = 012345678
INSERT INTO `Users`(pseudo , email , password,access) VALUES ("admin2","admin2@test.fr","$2a$10$zah2W8Y/jutBoUS/H3Zo.OABhYsC3qFrNesAifp6UK0KgvVj0pSwu",0);
-- Categories
INSERT INTO Categories(name) VALUES ("Patisserie");
INSERT INTO Categories(name) VALUES ("Plat");
INSERT INTO Categories(name) VALUES ("Soupe");
-- Posts
INSERT INTO Posts(title,imgUrl,description,date,user_id) 
VALUES (
    "Cookies",
    "/asset/image/Logo500G.png",
    "Ingredients 225g butter, softened 110g caster sugar 275g plain flour 1 tsp cinnamon or other spices (optional) 75g white or milk chocolate chips (optional)",
    "2015-08-10 19:21:00",
    (SELECT id FROM Users WHERE pseudo = "admin")
    );
INSERT INTO Posts(title,imgUrl,description,date,user_id) VALUES ("Lasagne","/asset/image/images.png","Lasagne de mamie","2018-05-20 19:21:00",(SELECT id FROM Users WHERE pseudo = "admin"));
INSERT INTO Posts(title,imgUrl,description,date,user_id) VALUES ("Soupe","/asset/image/images.png","Soupe de mamie","2020-05-20 19:21:00",(SELECT id FROM Users WHERE pseudo = "admin2"));
-- Categories link
INSERT INTO Categorized(post_id, category_id) VALUES (1,(SELECT id FROM Categories WHERE name = "Patisserie"));

-- Reaction
INSERT INTO Reactions(type,user_id,post_id) VALUES ("like",1,(SELECT id FROM Users WHERE pseudo = "admin")) ;
INSERT INTO Reactions(type,user_id,post_id) VALUES ("dislike",2,(SELECT id FROM Users WHERE pseudo = "admin")) ;