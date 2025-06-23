

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
	gender VARCHAR(10),
	country VARCHAR(25),
	religion VARCHAR(25),
	language VARCHAR(25),
	password VARCHAR(25)
	);

drop table profiles;

alter table enroll add create_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP;
alter table enroll drop matrimonyid;
alter table enroll add matrimonyid VARCHAR(25);
alter table enroll add looking VARCHAR(25);
select * from enrolls;
delete from enrolls where id > 0;

insert into enrolls (email, phone) values ('kandan@gmail.com','01123456789');
insert into profiles (email, phone) values ('kandan@gmail.com','01123456789');
delete from enroll where id > 0;
select * from profiles;
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



SELECT COUNT(*) FROM enroll WHERE email = 'dineshengg@gmail.com' OR phone = '9500008040';