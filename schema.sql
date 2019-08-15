CREATE TABLE character (
    character_id SERIAL PRIMARY KEY, 
    name VARCHAR(50) UNIQUE NOT NULL,
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