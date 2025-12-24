

create database kandanmatrimony;
USE kandanmatrimony;

CREATE TABLE enrolls (
	id SERIAL PRIMARY KEY,
	matrimonyid VARCHAR(25) unique not NULL,
	email VARCHAR(40) UNIQUE NOT NULL,
	phone VARCHAR(20) UNIQUE NOT null,
	looking VARCHAR(25),
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

create table profiles (
	id SERIAL primary key,
	matrimonyid VARCHAR(25) unique not NULL,
	firstname VARCHAR(40),
	secondname VARCHAR(40),
	email VARCHAR(40) UNIQUE NOT NULL,
	phone VARCHAR(20) UNIQUE NOT null,
	looking VARCHAR(25),
	DOB DATE,
	age INT,
	gender VARCHAR(10),
	country VARCHAR(25),
	religion VARCHAR(25),
	caste VARCHAR(25),
	state VARCHAR(40),
	city VARCHAR(40),
	language VARCHAR(25),
	password VARCHAR(25),
	hobbies VARCHAR(25),
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	status VARCHAR(15),
	subscription_type VARCHAR(15),
	verified BOOL,
	last_active TIMESTAMP,
	pref_age_min INT,
	pref_age_max INT,
	pref_religion VARCHAR(25),
	pref_caste VARCHAR(25),
    pref_country VARCHAR(40),
	pref_state VARCHAR(40),
	pref_city VARCHAR(40),
    pref_language VARCHAR(25),
	updated_at TIMESTAMP
	);


create table contactus (
	id SERIAL primary key,
	name VARCHAR(80),
	email VARCHAR(40),
	phone VARCHAR(20),
	category VARCHAR(40),
	message VARCHAR(256),
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

create table forgot (
	id SERIAL primary key,
	email VARCHAR(40),
	reset_at TIMESTAMP,
	times INT4,
	guid VARCHAR(36),
	matrimonyid VARCHAR(25)
	);
	
create table interests (
	id SERIAL primary key,
	sender_matid VARCHAR(25),
	receiver_matid VARCHAR(25),
	accepted BOOL,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

create table messages (
	id SERIAL primary key,
	sender_matid VARCHAR(25),
	receiver_matid VARCHAR(25),
	message VARCHAR(256),
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);


-- First, create the trigger function
CREATE OR REPLACE FUNCTION set_matrimonyid()
RETURNS TRIGGER AS $$
BEGIN
    NEW.matrimonyid := CONCAT('KAN', LPAD(NEW.id::text, 20, '0'));
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Then, create the trigger
CREATE TRIGGER before_insert_enroll
BEFORE INSERT ON enrolls
FOR EACH ROW
EXECUTE FUNCTION set_matrimonyid();
-- create trigger for the profiles
CREATE TRIGGER before_insert_profiles
BEFORE INSERT ON profiles
FOR EACH ROW
EXECUTE FUNCTION set_matrimonyid();



