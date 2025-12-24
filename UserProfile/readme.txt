//tables

CREATE TABLE profiles (
    id SERIAL PRIMARY KEY,
    matrimonyid VARCHAR(25) NOT NULL,
    firstname VARCHAR(50) NOT NULL,
    secondname VARCHAR(50) NOT NULL,
    email VARCHAR(40) UNIQUE NOT NULL,
    phone VARCHAR(20) UNIQUE NOT NULL,
    gender Gender,
    dob date NOT NULL,
    age INT NOT NULL,
    looking VARCHAR(20) NOT NULL,
    religion VARCHAR(20) NOT NULL,
    country VARCHAR(100) NOT NULL,
    language VARCHAR(25) NOT NULL,
    lastlogin timestamp NOT NULL,
    verified INT NOT NULL,
    photo_small VARCHAR(256) UNIQUE NULL, //store path where the photo got uploaded
    photo_approved BOOLEAN NOT NULL
)

CREATE TYPE Gender AS ENUM ('Male', 'Female');
CREATE TYPE MartialStatus AS ENUM ('Married', 'Un-married', 'Divorced', 'Awaiting-Divorce');

CREATE TABLE extprofiles (
    id INT FOREIGN KEY REFERENCES profiles(id) ON DELETE CASCADE,//delete if profile got deleted
    name VARCHAR(100) NOT NULL,//firstname and secondname
    age INT NOT NULL,
    height INT,
    weight INT,
    caste VARCHAR(40) NOT NULL,
    description TEXT,
    desc_approved BOOLEAN,
    gender Gender,
    color VARCHAR(10),
    photo1 VARCHAR(256) UNIQUE NULL,
    photo2 VARCHAR(256) UNIQUE NULL,
    photo3 VARCHAR(256) UNIQUE NULL,
    photo4 VARCHAR(256) UNIQUE NULL,
    photos_approved BOOLEAN,
    parentsno VARCHAR(20),
    alternate1 VARCHAR(20),
    father_name VARCHAR(25),
    mother_name VARCHAR(25),
    father_occupation VARCHAR(25),
    mother_occupation VARCHAR(25),
    family_info TEXT,
    famil_info_approved BOOLEAN,
    siblings_no INT,
    siblings_married_no INT,
    family_type VARCHAR(10),
    ancestors_origin VARCHAR(10),
    doj timestamp,
    marital_status MartialStatus NULL,
    disability_status BOOLEAN,
    disablility VARCHAR(256) NULL,
    state VARCHAR(40) NULL,
    city VARCHAR(40) NULL,
    location VARCHAR(100) nULL,
    sub_caste VARCHAR(40) NULL,
    education_level VARCHAR(10) NULL,
    education VARCHAR(25) NULL,
    working_status VARCHAR(25) NULL,
    industry VARCHAR(25) NULL,
    department VARCHAR(25) NULL,
    salary INT NULL,
    company VARCHAR(100),
    pan VARCHAR(12),
    aadhar VARCHAR(16),
    address VARCHAR(256),
    pincode INT,
    geolocation VARCHAR(100), //maps data like current location
    facebook VARCHAR(256),
    twitter VARCHAR(256),
    linkedin VARCHAR(256),
    others VARCHAR(256),
    hobby VARCHAR(256),
    reported INT,
    natchathiram VARCHAR(25),
    rasi VARCHAR(25),
);

CREATE TABLE preference (
    id INT FOREIGN KEY REFERENCES profiles(id) ON DELETE CASCADE,
    gender VARCHAR(10),
    age_min INT,
    age_max INT,
    religion VARCHAR(50),
    caste VARCHAR(50),
    sub_caste VARCHAR(40),
    language VARCHAR(50),
    state VARCHAR(50),
    country VARCHAR(50),
    working_status VARCHAR(20),
    salary_min INT,
    salary_max INT,
    marital_status VARCHAR(20),
    expectations TEXT
);
//membership tables


CREATE TABLE membership (
    id INT FOREIGN KEY REFERENCES profiles(id) ON DELETE CASCADE,
    membership_type VARCHAR(20) NOT NULL, -- e.g., 'free', 'trial', 'premium', 'paid'
    start_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    end_date TIMESTAMP,
    status BOOLEAN NOT NULL //default false based on expired or active.
);

//interest tables
CREATE TABLE interests (
    id SERIAL PRIMARY KEY,
    sender_id INT NOT NULL REFERENCES profiles(id) ON DELETE CASCADE,
    receiver_id INT NOT NULL REFERENCES profiles(id) ON DELETE CASCADE,
    status VARCHAR(20) NOT NULL DEFAULT 'pending', -- 'pending', 'accepted', 'rejected'
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);





////
chatgpt
Explanation of Changes
GORM Initialization:

Replaced sql.Open with gorm.Open using the PostgreSQL driver (gorm.io/driver/postgres).
Configured GORM with a logger for detailed query logging.
Connection Pooling:

Used SetMaxOpenConns to limit the maximum number of open connections to 20.
Used SetMaxIdleConns to limit the maximum number of idle connections to 15.
Used SetConnMaxLifetime to set the maximum lifetime of a connection to 30 minutes.
Ping Test:

Used sqlDB.Ping() to ensure the database connection is active.
Graceful Shutdown:

Updated CloseDB to close the GORM database connection properly.



////////////////////////////////////////
PGM:
Login flow - dev tested please test multiple times
Password reset flow
 - code done in backend not tested ********
 - python code for job worker needs to be written ********
20/12/25 
- Perfected login flow from home page
- DONE - Check the min age and max age for empty value 
- TODO - Enrolled but not put the full profile we have to send an email for incomplete profile
24/12/2025:
 - login flow completed
 - password reset flow completed
 - Check the min age and max age for empty value - done
 - 