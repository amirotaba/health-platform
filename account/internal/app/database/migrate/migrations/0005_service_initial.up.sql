CREATE TABLE IF NOT EXISTS service(
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(254) NULL UNIQUE,
    code VARCHAR(254) NOT NULL UNIQUE,
    path VARCHAR(254) NOT NULL,
    func VARCHAR(254) NOT NULL UNIQUE,
    method VARCHAR(254) NOT NULL,
    is_active bool,
    description VARCHAR(254),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    UNIQUE (path, method)
);

# CREATE TABLE IF NOT EXISTS actions(
#     id INT AUTO_INCREMENT PRIMARY KEY,
#     name VARCHAR(254) NOT NULL UNIQUE,
#     code VARCHAR(254) NOT NULL UNIQUE,
#     is_active bool,
#     description VARCHAR(254),
#     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ,
#     updated_at TIMESTAMP,
#     deleted_at TIMESTAMP
# );

CREATE TABLE IF NOT EXISTS permission_service(
    id INT AUTO_INCREMENT PRIMARY KEY,
    permission_id INT NOT NULL,
    service_id INT NOT NULL,
    description VARCHAR(254),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    FOREIGN KEY (service_id) REFERENCES service(id),
    FOREIGN KEY (permission_id) REFERENCES permission(id),
    UNIQUE (service_id, permission_id)
);
