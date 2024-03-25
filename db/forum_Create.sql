-- database: ./forum.sqlite

CREATE TABLE IF NOT EXISTS `Users` (
    username TEXT, 
    email TEXT, 
    password TEXT 
);


CREATE TABLE IF NOT EXISTS `Categories` (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT
);
CREATE TABLE IF NOT EXISTS `Posts` (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title VARCHAR(255) NOT NULL,
    imgUrl TEXT,
    description TEXT NOT NULL,
    date date NOT NULL,
    user_id INTEGER NOT NULL,

    FOREIGN KEY (user_id ) references Users(id)
);
CREATE TABLE IF NOT EXISTS `Comments` (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    message TEXT,
    date date NOT NULL,
    user_id INTEGER NOT NULL,
    post_id INTEGER NOT NULL,
    comment_id INTEGER,
    FOREIGN KEY (user_id ) references Users(id),
    FOREIGN KEY  (post_id ) references Posts(id),
    FOREIGN KEY (comment_id ) references Comments(id)
);
CREATE TABLE IF NOT EXISTS `Reactions` (
    type STRING NOT NULL,
    user_id INTEGER NOT NULL ,
    post_id INTEGER,
    comment_id INTEGER,
    
    FOREIGN KEY (user_id ) references Users(id),
    FOREIGN KEY  (post_id ) references Posts(id),
    FOREIGN KEY  (comment_id ) references Comments(id),
    PRIMARY KEY (user_id, post_id,comment_id)
);
CREATE TABLE IF NOT EXISTS `Categorized` (
    category_id INTEGER NOT NULL,
    post_id INTEGER NOT NULL,

    FOREIGN KEY (category_id ) references Categories(id),
    FOREIGN KEY (post_id ) references Posts(id),
    PRIMARY KEY (category_id, post_id)
);
