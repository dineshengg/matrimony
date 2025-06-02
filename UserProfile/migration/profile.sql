CREATE TABLE profiles (
    id SERIAL PRIMARY KEY,
    //TODO - matrimony id add sql INCREMENT with some prefix
    matrimonyid VARCHAR(25) NOT NULL,
    firstname VARCHAR(50) NOT NULL,
    secondname VARCHAR(50) NOT NULL,
    email VARCHAR(40) UNIQUE NOT NULL,
    password VARCHAR(12) UNIQUE NOT NULL,
    email_verified BOOLEAN NOT NULL,
    phone VARCHAR(20) UNIQUE NOT NULL,
    phone_verified BOOLEAN NOT NULL,
    gender Gender,
    dob date NOT NULL,
    looking VARCHAR(20) NOT NULL,
    religion VARCHAR(20) NOT NULL,
    country VARCHAR(100) NOT NULL,
    language VARCHAR(25) NOT NULL,
    createdat timestamp NOT NULL,
    verified INT NOT NULL
)

CREATE TYPE Gender AS ENUM ('Male', 'Female');
