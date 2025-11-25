CREATE TABLE IF NOT EXISTS staff (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    hospital VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS patient (
    id SERIAL PRIMARY KEY,
    national_id VARCHAR(50),
    passport_id VARCHAR(50),
    first_name VARCHAR(255) NOT NULL,
    middle_name VARCHAR(255),
    last_name VARCHAR(255) NOT NULL,
    date_of_birth DATE NOT NULL,
    phone_number VARCHAR(50),
    email VARCHAR(255),
    gender VARCHAR(10) NOT NULL,
    patient_hn VARCHAR(50) NOT NULL,
    hospital VARCHAR(255) NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(patient_hn, hospital)
);

CREATE INDEX IF NOT EXISTS idx_staff_username ON staff(username);
CREATE INDEX IF NOT EXISTS idx_staff_hospital ON staff(hospital);
CREATE INDEX IF NOT EXISTS idx_patient_hn ON patient(patient_hn);
CREATE INDEX IF NOT EXISTS idx_patient_hospital ON patient(hospital);

