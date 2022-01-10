CREATE TYPE sex AS ENUM ('F', 'M');

CREATE TABLE IF NOT EXISTS players(
    id SERIAL,
    username VARCHAR(16) NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    password_salt BYTEA NOT NULL,
    sso_token TEXT DEFAULT NULL,
    sex SEX NOT NULL DEFAULT 'F',
    figure TEXT NOT NULL DEFAULT '1000118001270012900121001',
    pool_figure TEXT NOT NULL DEFAULT '',
    film INT NOT NULL DEFAULT 0,
    credits INT NOT NULL DEFAULT 100,
    tickets INT NOT NULL DEFAULT 0,
    motto TEXT NOT NULL DEFAULT 'Project habbgo.',
    console_motto TEXT NOT NULL DEFAULT 'habbgo rocks!',
    birthday DATE NOT NULL,
    email TEXT NOT NULL,
    sound_enabled BOOL NOT NULL DEFAULT true,
    created_on TIMESTAMP NOT NULL,
    last_online TIMESTAMP NOT NULL DEFAULT current_timestamp,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS badges(
    id   SERIAL,
    code VARCHAR(3) UNIQUE NOT NULL,
    display BOOL NOT NULL DEFAULT false,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS player_badges (
    player_id INT NOT NULL,
    badge_id INT NOT NULL,
    FOREIGN KEY (player_id) REFERENCES Players(id) ON DELETE CASCADE,
    FOREIGN KEY (badge_id) REFERENCES Badges(id),
    PRIMARY KEY (player_id, badge_id)
);

CREATE TABLE IF NOT EXISTS room_categories (
    id SERIAL,
    parent_id INT NOT NULL,
    order_id INT NOT NULL,
    name TEXT NOT NULL,
    is_node BOOL NOT NULL DEFAULT false,
    is_public BOOL NOT NULL DEFAULT false,
    is_trading BOOL NOT NULL DEFAULT false,
    min_rank_access INT NOT NULL DEFAULT 1,
    min_rank_setflatcat INT NOT NULL DEFAULT 1,
    PRIMARY KEY (id)
);

INSERT INTO room_categories (id, order_id, parent_id, is_node, name, is_public, is_trading, min_rank_access, min_rank_setflatcat)
VALUES (2, 0, 0, false, 'No category', false, false, 1, 1),
       (3, 0, 0, true, 'Public Rooms', true, false, 1, 6),
       (4, 0, 0, true, 'Guest Rooms', false, false, 1, 6),
       (5, 0, 3, false, 'Entertainment', true, false, 1, 6),
       (6, 0, 3, false, 'Restaurants and Cafes', true, false, 1, 6),
       (7, 0, 3, false, 'Lounges and Clubs', true, false, 1, 6),
       (8, 0, 3, false, 'Club-only Spaces', true, false, 1, 6),
       (9, 0, 3, false, 'Parks and Gardens', true, false, 1, 6),
       (10, 0, 3, false, 'Swimming Pools', true, false, 1, 6),
       (11, 0, 3, false, 'The Lobbies', true, false, 1, 6),
       (12, -1, 3, false, 'The Hallways', true, false, 1, 6),
       (13, 0, 3, false, 'Games', true, false, 1, 6),
       (101, 0, 4, false, 'Staff HQ', false, true, 4, 5),
       (112, 0, 4, false, 'Restaurant, Bar & Night Club Rooms', false, false, 1, 1),
       (113, 0, 4, false, 'Trade floor', false, true, 1, 1),
       (114, 0, 4, false, 'Chill, Chat & Discussion Rooms', false, false, 1, 1),
       (115, 0, 4, false, 'Hair Salons & Modelling Rooms', false, false, 1, 1),
       (116, 0, 4, false, 'Maze & Theme Park Rooms', false, false, 1, 1),
       (117, 0, 4, false, 'Gaming & Race Rooms', false, false, 1, 1),
       (118, 0, 4, false, 'Help Centre Rooms', false, false, 1, 1),
       (120, 0, 4, false, 'Miscellaneous', false, false, 1, 1);
