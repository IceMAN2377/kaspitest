CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    phone VARCHAR(20) NOT NULL UNIQUE,
    iin VARCHAR(12) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL
);