CREATE TABLE IF NOT EXISTS token(
    id INT AUTO_INCREMENT PRIMARY KEY,
    uuid VARCHAR(254),
    account_id INT,
    auth_token TEXT,
    refresh_token TEXT,
    type VARCHAR(254),
    expire_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) ENGINE=INNODB;
