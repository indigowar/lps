--
-- This file contains overall schema of the project
--

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS positions(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v1(),
    title VARCHAR(128) UNIQUE NOT NULL,
    max_per_department INTEGER
);

CREATE TABLE IF NOT EXISTS departments(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v1(),
    name VARCHAR(256) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS staff(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v1(),
    surname VARCHAR(64) NOT NULL,
    name VARCHAR(64) NOT NULL,
    patronymic VARCHAR(64),
    phone_number VARCHAR(15) UNIQUE NOT NULL,
    position UUID REFERENCES positions(id) NOT NULL,
    department UUID REFERENCES departments(id) NOT NULL
);

CREATE TABLE IF NOT EXISTS incidents(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v1(),
    employee UUID REFERENCES staff(id) NOT NULL,
    description TEXT NOT NULL,
    happenning_date TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS professional_developments(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v1(),
    employee UUID REFERENCES staff(id) NOT NULL,
    title VARCHAR(255) NOT NULL,
    starting_date DATE NOT NULL,
    ending_date DATE NOT NULL
);

CREATE TABLE IF NOT EXISTS accounts(
    login VARCHAR(255) PRIMARY KEY,
    password VARCHAR(1024),
    activated BOOLEAN NOT NULL DEFAULT FALSE,
    employee UUID REFERENCES staff(id)
        ON UPDATE CASCADE
        ON DELETE CASCADE
        NOT NULL UNIQUE
);
