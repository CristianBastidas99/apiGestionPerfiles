CREATE TABLE IF NOT EXISTS user_profiles (
    UserID INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    URL VARCHAR(255),
    Nickname VARCHAR(100),
    ContactPublic TINYINT(1),
    Address VARCHAR(255),
    Biography TEXT,
    Organization VARCHAR(100),
    Country VARCHAR(100),
    SocialLinks TEXT
);
