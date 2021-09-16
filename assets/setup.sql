CREATE DATABASE habbgo;

CREATE TABLE Player {
    ID INT NOT NULL UNIQUE,
    Username VARCHAR(16) NOT NULL UNIQUE,
    Sex CHAR(1) NOT NULL,
    Figure TEXT NOT NULL,
    PoolFigure TEXT,
    Film INT DEFAULT 0,
    Credits INT DEFAULT 0,
    Tickets INT DEFAULT 0,
    Motto TEXT DEFAULT 'Project HabbGo.',
    ConsoleMotto TEXT DEFAULT 'HabbGo Rocks!',
    DisplayBadge BOOL DEFAULT false,
    CurrentBadge INT,
    SoundEnabled BOOL DEFAULT true,
    CreatedOn DATETIME NOT NULL,
    LastOnline DATETIME,
    FOREIGN KEY (CurrentBadge) REFERENCES Badge(ID)
    PRIMARY KEY (ID)
};

CREATE TABLE Badge {
    ID INT NOT NULL UNIQUE,
    Code VARCHAR(6) UNIQUE,
    PRIMARY KEY (ID)
};

CREATE TABLE PlayerBadges {
    PlayerID INT NOT NULL,
    Badge INT,
    FOREIGN KEY (PlayerID) REFERENCES Player(ID)
    FOREIGN KEY (Badge) REFERENCES Badge(ID)
    PRIMARY KEY (PlayerID)
};