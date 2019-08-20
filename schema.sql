CREATE TABLE character (
    id SERIAL, 
    name VARCHAR(50) UNIQUE NOT NULL,
    display_name VARCHAR(50) UNIQUE NOT NULL,
    icon VARCHAR(255) UNIQUE,
    PRIMARY KEY(id)
);

CREATE TABLE lesson (
    id SERIAL, 
    character_id INTEGER REFERENCES character(id) NOT NULL, 
    name VARCHAR(50) UNIQUE NOT NULL,
    number INTEGER NOT NULL,
    gif VARCHAR(255),
    description TEXT,
    learning_time_seconds INTEGER,
    training_time_seconds INTEGER,
    test_time_seconds INTEGER,
    PRIMARY KEY(id)
);