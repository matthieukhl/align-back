-- Create the Pilates Management Database
CREATE DATABASE IF NOT EXISTS pilates_management;
USE pilates_management;

-- Enable foreign key constraints
SET FOREIGN_KEY_CHECKS = 1;

-- Clients Table
CREATE TABLE IF NOT EXISTS clients (
    id VARCHAR(36) PRIMARY KEY,
    full_name VARCHAR(255) NOT NULL,
    firstname VARCHAR(100) NOT NULL,
    lastname VARCHAR(100) NOT NULL,
    phone VARCHAR(20),
    email VARCHAR(100) UNIQUE,
    street_number VARCHAR(20),
    street_name VARCHAR(255),
    city VARCHAR(100),
    zip_code VARCHAR(20),
    country VARCHAR(100),
    credits INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Packages Table
CREATE TABLE IF NOT EXISTS packages (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    number_of_sessions INT NOT NULL,
    type ENUM('GROUP', 'PRIVATE') NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Classes Table
CREATE TABLE IF NOT EXISTS classes (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    location ENUM('CLAIRVIVRE', 'CUBJAC') NOT NULL,
    type ENUM('GROUP', 'PRIVATE') NOT NULL,
    equipment VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Schedule Table
CREATE TABLE IF NOT EXISTS schedule (
    id VARCHAR(36) PRIMARY KEY,
    class_id VARCHAR(36) NOT NULL,
    capacity INT NOT NULL DEFAULT 10,
    class_datetime DATETIME NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (class_id) REFERENCES classes(id) ON DELETE CASCADE
);

-- Appointments Table (renamed from appointment for consistency)
CREATE TABLE IF NOT EXISTS appointments (
    id VARCHAR(36) PRIMARY KEY,
    schedule_id VARCHAR(36) NOT NULL,
    client_id VARCHAR(36) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (schedule_id) REFERENCES schedule(id) ON DELETE CASCADE,
    FOREIGN KEY (client_id) REFERENCES clients(id) ON DELETE CASCADE,
    UNIQUE KEY unique_appointment (schedule_id, client_id)
);

-- Billing Table
CREATE TABLE IF NOT EXISTS billings (
    id VARCHAR(36) PRIMARY KEY,
    client_id VARCHAR(36) NOT NULL,
    package_id VARCHAR(36) NOT NULL,
    amount INT NOT NULL DEFAULT 1,
    price DECIMAL(10, 2) NOT NULL,
    credits INT NOT NULL,
    payment_date DATETIME DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (client_id) REFERENCES clients(id) ON DELETE CASCADE,
    FOREIGN KEY (package_id) REFERENCES packages(id) ON DELETE CASCADE
);

-- Create indices for performance
CREATE INDEX idx_clients_email ON clients(email);
CREATE INDEX idx_clients_name ON clients(lastname, firstname);
CREATE INDEX idx_schedule_datetime ON schedule(class_datetime);
CREATE INDEX idx_appointments_client ON appointments(client_id);
CREATE INDEX idx_appointments_schedule ON appointments(schedule_id);
CREATE INDEX idx_billings_client ON billings(client_id);

-- Insert some sample data
INSERT INTO clients (id, full_name, firstname, lastname, phone, email) VALUES
(UUID(), 'Jane Smith', 'Jane', 'Smith', '+33123456789', 'jane.smith@example.com'),
(UUID(), 'John Doe', 'John', 'Doe', '+33987654321', 'john.doe@example.com');

INSERT INTO packages (id, name, number_of_sessions, type, price) VALUES
(UUID(), '10 Group Sessions', 10, 'GROUP', 150.00),
(UUID(), '5 Private Sessions', 5, 'PRIVATE', 250.00);

INSERT INTO classes (id, name, location, type, equipment) VALUES
(UUID(), 'Morning Pilates', 'CLAIRVIVRE', 'GROUP', 'Mat, resistance bands'),
(UUID(), 'Reformer Session', 'CUBJAC', 'PRIVATE', 'Reformer machine');