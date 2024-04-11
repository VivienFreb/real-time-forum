-- database: ./forum.sqlite

CREATE TABLE IF NOT EXISTS `Users` (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT, 
    email TEXT, 
    password TEXT,
    status TEXT 
);

CREATE TABLE IF NOT EXISTS `Posts` (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    user_id INTEGER NOT NULL,

    FOREIGN KEY (user_id ) references Users(id)
);

CREATE TABLE IF NOT EXISTS `Disscussions`(
    speaker TEXT NOT NULL,
    listener TEXT NOT NULL,
    content TEXT NOT NULL
)

-- CREATE TABLE IF NOT EXISTS `Comments` (
--     id INTEGER PRIMARY KEY AUTOINCREMENT,
--     message TEXT,
--     date date NOT NULL,
--     user_id INTEGER NOT NULL,
--     post_id INTEGER NOT NULL,
--     comment_id INTEGER,
--     FOREIGN KEY (user_id ) references Users(id),
--     FOREIGN KEY  (post_id ) references Posts(id),
--     FOREIGN KEY (comment_id ) references Comments(id)
-- );
