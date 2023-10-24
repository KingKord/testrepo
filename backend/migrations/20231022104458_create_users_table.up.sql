-- 001_create_users_table.up.sql
CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       name VARCHAR(255) NOT NULL,
                       surname VARCHAR(255) NOT NULL,
                       patronymic VARCHAR(255),
                       agify VARCHAR(255),
                       genderize VARCHAR(255),
                       nationalize VARCHAR(255)
);
