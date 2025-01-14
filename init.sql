-- create the users table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    surname VARCHAR(50) NOT NULL,
    email VARCHAR(100) NOT NULL UNIQUE
);

-- fill it with data
INSERT INTO users (name, surname, email) VALUES
('John', 'Doe', 'john.doe@example.com'),
('Guido', 'Enrique', 'guidoenr4@gmail.com'),
('Bel', 'Bianchi', 'belubelu@gmail.com'),
('Kevin', 'Nurray', 'kevinnurray@proton.email'),
('Elliot', 'Anderson', 'elliotanderson@hotmail.com'),
('Albert', 'Brown', 'albet_brown_02@hotmail.com');


