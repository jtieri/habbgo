CREATE DATABASE IF NOT EXISTS `habbgo`;

CREATE TABLE IF NOT EXISTS `Players` (
    `ID` INT PRIMARY KEY AUTO_INCREMENT,
    `Username` VARCHAR(16) NOT NULL UNIQUE,
    `PasswordHash` TEXT NOT NULL,
    `PasswordSalt` VARBINARY NOT NULL,
    `SSOToken` TEXT DEFAULT NULL,
    `Sex` ENUM('M','F') NOT NULL DEFAULT 'F',
    `Figure` TEXT NOT NULL DEFAULT '1000118001270012900121001',
    `PoolFigure` TEXT,
    `Film` INT DEFAULT 0,
    `Credits` INT DEFAULT 100,
    `Tickets` INT DEFAULT 0,
    `Motto` TEXT DEFAULT 'Project HabbGo.',
    `ConsoleMotto` TEXT DEFAULT 'HabbGo Rocks!',
    `DisplayBadge` BOOL NOT NULL DEFAULT true,
    `CurrentBadge` INT,
    `Birthday` DATE NOT NULL,
    `Email` TEXT NOT NULL,
    `SoundEnabled` BOOL NOT NULL DEFAULT true,
    `CreatedOn` DATETIME NOT NULL,
    `LastOnline` DATETIME,
    FOREIGN KEY (CurrentBadge) REFERENCES Badges (ID)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `Badges` (
    `ID` INT PRIMARY KEY AUTO_INCREMENT,
    `Code` VARCHAR(3) UNIQUE NOT NULL
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `PlayerBadges` (
    `PlayerID` INT,
    `Badge` INT,
    FOREIGN KEY (PlayerID) REFERENCES Players (ID) ON DELETE CASCADE,
    FOREIGN KEY (Badge) REFERENCES Badges (ID)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;