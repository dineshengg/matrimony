

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

alter table enroll add create_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP;
alter table enroll drop matrimonyid;
alter table enroll add matrimonyid VARCHAR(25);
alter table enroll add looking VARCHAR(25);
select * from enroll;

insert into enroll (email, phone) values ('kandan1@gmail.com','01123456789');
delete from enroll where id > 0;
-- First, create the trigger function
CREATE OR REPLACE FUNCTION set_matrimonyid()
RETURNS TRIGGER AS $$
BEGIN
    NEW.matrimonyid := CONCAT('KAN', LPAD(NEW.id::text, 6, '0'));
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Then, create the trigger
CREATE TRIGGER before_insert_enroll
BEFORE INSERT ON enroll
FOR EACH ROW
EXECUTE FUNCTION set_matrimonyid();


SELECT COUNT(*) FROM enroll WHERE email = 'dineshengg@gmail.com' OR phone = '9500008040';