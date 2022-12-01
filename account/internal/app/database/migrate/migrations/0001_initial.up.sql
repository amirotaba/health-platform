CREATE TABLE IF NOT EXISTS permission(
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(254) NOT NULL UNIQUE,
    is_active bool,
    description VARCHAR(254),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS account_type(
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(254) NOT NULL UNIQUE,
    is_active bool,
    description VARCHAR(254),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS membership(
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(254) NOT NULL UNIQUE,
    is_active bool,
    description VARCHAR(254),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS role(
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(254) NOT NULL UNIQUE,
    is_active bool,
    description VARCHAR(254),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS account(
    id INT AUTO_INCREMENT PRIMARY KEY,
    uuid VARCHAR(254) NOT NULL,
    first_name VARCHAR(254) NOT NULL,
    last_name VARCHAR(254) NOT NULL,
    display_name VARCHAR(254),
    password VARCHAR(254),
    email VARCHAR(254),
    phone_number VARCHAR(254) UNIQUE,
    address VARCHAR(254),
    image_url VARCHAR(254),
    type_id INT,
    role_id INT,
    is_active BOOLEAN,
    expire_at TIMESTAMP,
    last_login TIMESTAMP,
    description VARCHAR(254),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    CONSTRAINT fk_account_type FOREIGN KEY (type_id) REFERENCES account_type(id) ON DELETE SET NULL,
    CONSTRAINT fk_account_role FOREIGN KEY (role_id) REFERENCES role(id) ON DELETE SET NULL
) ENGINE=INNODB;

CREATE TABLE IF NOT EXISTS channel (
    id INT AUTO_INCREMENT PRIMARY KEY,
    uuid VARCHAR(254) NOT NULL,
    name VARCHAR(254) NOT NULL UNIQUE,
    display_name VARCHAR(254),
    password VARCHAR(254),
    email VARCHAR(254),
    image_url VARCHAR(254),
    current_balance DECIMAL DEFAULT 0.0,
    membership_type_id INT,
    is_active BOOLEAN,
    expire_at TIMESTAMP,
    owner_phone_number VARCHAR(254),
    description VARCHAR(254),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    CONSTRAINT fk_channel_membership_type FOREIGN KEY (membership_type_id) REFERENCES account_type(id) ON DELETE SET NULL
) ENGINE=INNODB;

CREATE TABLE IF NOT EXISTS channel_rule(
    id INT AUTO_INCREMENT PRIMARY KEY,
    channel_id INT NOT NULL,
    tag_id INT NOT NULL,
    price DECIMAL NOT NULL DEFAULT 0.0,
    is_active BOOLEAN,
    destination VARCHAR(254),
    description VARCHAR(254),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    UNIQUE (channel_id, tag_id),
    CONSTRAINT fk_channel_rule_channel_id FOREIGN KEY (channel_id) REFERENCES channel(id)
);

CREATE TABLE IF NOT EXISTS channel_admin(
    id INT AUTO_INCREMENT PRIMARY KEY,
    account_id INT NOT NULL,
    role_id INT NOT NULL,
    channel_id INT NOT NULL,
    description VARCHAR(254),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    FOREIGN KEY (role_id) REFERENCES role(id),
    FOREIGN KEY (account_id) REFERENCES account(id),
    FOREIGN KEY (channel_id) REFERENCES channel(id)
);

CREATE TABLE IF NOT EXISTS role_permission(
    id INT AUTO_INCREMENT PRIMARY KEY,
    role_id INT NOT NULL,
    permission_id INT NOT NULL,
    description VARCHAR(254),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    FOREIGN KEY (role_id) REFERENCES role(id),
    FOREIGN KEY (permission_id) REFERENCES permission(id),
    UNIQUE (role_id, permission_id)
);
