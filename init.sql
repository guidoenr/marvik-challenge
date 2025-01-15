-- create the users table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    surname VARCHAR(50) NOT NULL,
    email VARCHAR(100) NOT NULL UNIQUE
);

-- create the organizations table
CREATE TABLE IF NOT EXISTS organizations (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL
);

-- create the user_organizations (join) table
CREATE TABLE IF NOT EXISTS user_organizations (
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    organization_id INT REFERENCES organizations(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, organization_id)
);


-- fill users table with sample data
INSERT INTO users (name, surname, email) VALUES
('John', 'Doe', 'john.doe@example.com'),
('Guido', 'Enrique', 'guidoenr4@gmail.com'),
('Bel', 'Bianchi', 'belubelu@gmail.com'),
('Kevin', 'Nurray', 'kevinnurray@proton.email'),
('Elliot', 'Anderson', 'elliotanderson@hotmail.com'),
('Susana', 'Gimenez', 'susanita@yahoo.com.ar'),
('Albert', 'Brown', 'albet_brown_02@hotmail.com');

-- fill organizations table with sample data
INSERT INTO organizations (name) VALUES
('Veritone'),
('Slingr'),
('Marvik'),
('Google'),
('Amazon'),
('Hackerone');

-- assign users to organizations (many-to-many relationship)
INSERT INTO user_organizations (user_id, organization_id) VALUES
(1, 1), -- John Doe belongs to Veritone
(2, 2), -- Guido Enrique belongs to Slingr
(2, 1), -- and go on ...
(2, 4), 
(3, 3), 
(4, 4),
(4, 5), 
(5, 1), 
(6, 2),
(7, 1),
(7, 2);

