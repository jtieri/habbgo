CREATE DATABASE IF NOT EXISTS `habbgo`;

CREATE TABLE IF NOT EXISTS `Player` (
    `ID` INT,
    `Username` VARCHAR(16) NOT NULL UNIQUE,
    `Password` VARCHAR(10) NOT NULL,
    `SSOToken` TEXT DEFAULT NULL,
    `Sex` ENUM('M','F') NOT NULL DEFAULT 'F',
    `Figure` TEXT NOT NULL DEFAULT '1000118001270012900121001',
    `PoolFigure` TEXT,
    `Film` INT DEFAULT 0,
    `Credits` INT DEFAULT 100,
    `Tickets` INT DEFAULT 0,
    `Motto` TEXT DEFAULT 'Project HabbGo.',
    `ConsoleMotto` TEXT DEFAULT 'HabbGo Rocks!',
    `DisplayBadge` BOOL DEFAULT false,
    `CurrentBadge` INT,
    `SoundEnabled` BOOL DEFAULT true,
    `CreatedOn` DATETIME NOT NULL,
    `LastOnline` DATETIME,
    PRIMARY KEY (ID),
    FOREIGN KEY (CurrentBadge) REFERENCES Badge(ID)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `Badge` (
    `ID` INT PRIMARY KEY AUTO_INCREMENT,
    `Code` VARCHAR(3) UNIQUE NOT NULL
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `PlayerBadges` (
    `PlayerID` INT,
    `Badge` INT,
    FOREIGN KEY (PlayerID) REFERENCES Player(ID),
    FOREIGN KEY (Badge) REFERENCES Badge(ID)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;