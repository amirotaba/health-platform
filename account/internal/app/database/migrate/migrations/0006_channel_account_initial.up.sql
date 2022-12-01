CREATE TABLE IF NOT EXISTS channel_account(
    id INT AUTO_INCREMENT PRIMARY KEY,
    channel_id INT NOT NULL,
    account_id INT NOT NULL,
    role_id INT NOT NULL,
    description VARCHAR(254),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    FOREIGN KEY (account_id) REFERENCES account(id),
    FOREIGN KEY (channel_id) REFERENCES channel(id),
    FOREIGN KEY (role_id) REFERENCES role(id),
    UNIQUE (account_id, channel_id)
);

CREATE TABLE IF NOT EXISTS account_role(
    id INT AUTO_INCREMENT PRIMARY KEY,
    account_id INT NOT NULL,
    role_id INT NOT NULL,
    description VARCHAR(254),
    is_active VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    FOREIGN KEY (account_id) REFERENCES account(id),
    FOREIGN KEY (role_id) REFERENCES role(id),
    UNIQUE (account_id, role_id)
);