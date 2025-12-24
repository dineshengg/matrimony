

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

select * from enrolls;
insert into enrolls (email, phone) values ('kandan@gmail.com','01123456789');
delete from enrolls where id > 0;
SELECT COUNT(*) FROM enroll WHERE email = 'dineshengg@gmail.com' OR phone = '9500008040';

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

------------

insert into profiles (firstname,secondname,email,phone,looking,DOB,gender,country,religion,language,password,status,subscription_type,verified) values('ganesh', 'murthy','ganeshmurthy1@gmail.com', '951000008040','son','2025-12-09', 'Male', 'India','Hindu', 'Tamil', '$2a$10$rR2v5vg1cbKm8BOPpz1kseiU4bvK9T/FSpG.R7mqLo6SqQJN9HyNa', 'active','trial',false);
delete from profiles where matrimonyid='KAN00000000000000000028';

insert into profiles (email, phone) values ('kandan@gmail.com','01123456789');
drop table profiles;
select * from profiles;
alter table profile add status VARCHAR(15);
alter table profiles add pref_age_min INT;
alter table profiles add pref_age_max INT;
alter table profiles add pref_religion VARCHAR(25);
alter table profiles add caste VARCHAR(25);
alter table profiles add pref_caste VARCHAR(25);
alter table profiles add pref_state VARCHAR(40);
alter table profiles add hobbies VARCHAR(25);
alter table profiles add state VARCHAR(40);
alter table profiles add city VARCHAR(40);
alter table profiles add updated_at TIMESTAMP;
alter table profiles add pref_city VARCHAR(40);
alter table profiles add age INT;
alter table profiles add pref_country VARCHAR(40);
alter table profiles add pref_language VARCHAR(25);

ALTER TABLE profiles DROP COLUMN pref_state;
ALTER TABLE profiles DROP COLUMN location;
ALTER TABLE profiles DROP COLUMN state;
---

create table counters (
	matid_category VARCHAR(24),
	counter integer
	
);

select * from counters;
drop table counters;

---
create table globalcounters (
	category VARCHAR(24) unique not null,
	counter integer
);
select * from globalcounters;
insert into globalcounters (category, counter) values ('mathtt', 1) on conflict (category) do update set category=EXCLUDED.category, counter=EXCLUDED.counter;
select counter from globalcounters order by counter desc limit 1;
drop table globalcounters;

----
UPDATE forgot SET  guid = 'asdfghjkl', times =  1 WHERE email = 'techlearningbox@gmail.com';
SELECT email, guid FROM forgot WHERE times > 0;
alter table forgot add matrimonyid VARCHAR(25);
---------
alter table enroll add create_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP;
alter table enroll drop matrimonyid;
alter table enroll add matrimonyid VARCHAR(25);
alter table enroll add looking VARCHAR(25);

-------
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
	accepted VARCHAR(15),
    email_sent BOOL NOT null default false,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
 drop table interests;


create table messages (
	id SERIAL primary key,
	sender_matid VARCHAR(25),
	receiver_matid VARCHAR(25),
	message VARCHAR(256),
	email_sent BOOL NOT null default false,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
drop table messages;
	
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






