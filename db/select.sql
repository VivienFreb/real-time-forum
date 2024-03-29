-- database: ./forum.sqlite

-- User 
SELECT u.pseudo, p.description, p.name
FROM Posts p
INNER JOIN Users u ON p.'user_id' = u.'id';
-- User by Id
SELECT * 
FROM Users 
WHERE id = 1;
-- User pseudo by Id
SELECT pseudo 
FROM Users 
WHERE id = 2;

-- All Posts ordered by date and description
SELECT * 
FROM Posts 
ORDER BY date DESC;

-- Category by Id
SELECT * 
FROM Categories 
WHERE id = 1;

-- Categories to a post
SELECT * 
FROM Categorized 
WHERE post_id;

SELECT COUNT(*) 
FROM Users 
WHERE pseudo = "admin" AND mdp = "1234";

SELECT COUNT(*) AS NBlike 
FROM Reactions
INNER JOIN Posts ON id = post_id 
WHERE post_id = 1 AND type = "like";

SELECT COUNT(*) AS NBdislike 
FROM Reactions
INNER JOIN Posts ON id = post_id 
WHERE post_id = 1 AND type = "dislike";


SELECT
COUNT(CASE WHEN type = "like" THEN 1 END) AS NBlike,
COUNT(CASE WHEN type = "dislike" THEN 1 END) AS NBdislike 
FROM Reactions
INNER JOIN Posts ON id = post_id 
WHERE post_id = 1;

-- Post 
SELECT * 
FROM Posts 
WHERE user_id = 1 AND name = 'Cookies' LIMIT 5;

SELECT
COUNT(id) AS NBid
FROM Posts;

SELECT u.pseudo, c.message, c.date FROM Comments c
INNER JOIN Posts p ON c.post_id = p.id
INNER JOIN Users u ON c.user_id = u.id;
-- Where c.post_id = 2;

SELECT * FROM Comments;

-- Fetch Posts Table rows depending on Category Table name
SELECT p.id, p.title, p.imgUrl, p.description, p.date, p.user_id, u.pseudo FROM Posts p
INNER JOIN Users u ON p.user_id = u.id
RIGHT JOIN Categorized c ON c.post_id = p.id
INNER JOIN Categories ca ON ca.id = c.category_id 
WHERE ca.name = 'Patisserie';

-- Fetch Posts Table rows depending on user.pseudo
SELECT p.id, p.title, p.imgUrl, p.description, p.date, p.user_id, u.pseudo  FROM Posts p
INNER JOIN Users u  ON p.user_id = u.id
WHERE u.pseudo = 'admin';

-- Fetch one row of Posts Table with Reaction like and dislike

SELECT DISTINCT p.id, p.title, p.imgUrl, p.description, p.date, u.pseudo,
(SELECT COUNT(CASE WHEN type = "like" THEN 1 END) 
FROM Reactions
INNER JOIN Posts ON id = post_id 
WHERE post_id = p.id) AS NBlike,
(SELECT COUNT(CASE WHEN type = "dislike" THEN 1 END) 
FROM Reactions
INNER JOIN Posts ON id = post_id 
WHERE post_id = p.id) AS NBdislike
FROM Posts p
INNER JOIN Users u  ON p.user_id = u.id
INNER JOIN Reactions r ON p.id = r.post_id
ORDER BY p.date ASC, p.title DESC ;
-- LEFT JOIN Categorized c ON c.post_id = p.id
-- INNER JOIN Categories ca ON ca.id = c.category_id 

SELECT DISTINCT p.id, p.Title, p.imgUrl, p.Description, p.date, u.pseudo, 
(SELECT COUNT(CASE WHEN type = "like" THEN 1 END)
FROM Reactions 
INNER JOIN Posts ON id = post_id 
WHERE post_id = p.id) AS NBlike, 
(SELECT COUNT(CASE WHEN type = "dislike" THEN 1 END)
FROM Reactions 
INNER JOIN Posts ON id = post_id 
WHERE post_id = p.id) AS NBdislike
FROM Posts p
INNER JOIN Users u ON p.user_id = u.id 
INNER JOIN Reactions r ON p.id = r.post_id 
ORDER BY p.date ASC, p.title DESC ;

     
SELECT DISTINCT p.id, p.Title, p.imgUrl, p.Description, p.date, u.pseudo, 
(SELECT COUNT(CASE WHEN type = "like" THEN 1 END) 
FROM Reactions 
INNER JOIN Posts ON id = post_id 
WHERE post_id = p.id) AS NBlike, 
(SELECT COUNT(CASE WHEN type = "dislike" THEN 1 END) 
FROM Reactions 
INNER JOIN Posts ON id = post_id 
WHERE post_id = p.id) AS NBdislike, 
(SELECT COUNT(*) 
FROM Comments 
WHERE post_id = p.id) AS NBcomment 
FROM Posts p 
INNER JOIN Users u ON p.user_id = u.id 
INNER JOIN Reactions r ON p.id = r.post_id 
LEFT JOIN Categorized c ON c.post_id = p.id 
LEFT JOIN Categories ca ON ca.id = c.category_id 
WHERE ca.name = "Patisserie"
ORDER BY p.Title;