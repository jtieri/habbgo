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
    FOREIGN KEY (badge_id) REFERENCES Badges(id)
);