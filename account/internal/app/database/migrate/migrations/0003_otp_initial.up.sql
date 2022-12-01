CREATE TABLE IF NOT EXISTS account_otp(
    id INT AUTO_INCREMENT PRIMARY KEY,
    account_id INT,
    phone_number VARCHAR(254),
    otp VARCHAR(254),
    expire_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_account_id FOREIGN KEY (account_id) REFERENCES account(id) ON DELETE SET NULL
) ENGINE=INNODB;
