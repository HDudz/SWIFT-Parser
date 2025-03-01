CREATE TABLE swiftTable (
                             id INT AUTO_INCREMENT PRIMARY KEY,
                             country_iso2 CHAR(2) NOT NULL,
                             code VARCHAR(11) NOT NULL UNIQUE,
                             code_type VARCHAR(50) NOT NULL,
                             bank_name VARCHAR(255) NOT NULL,
                             address VARCHAR(255) NOT NULL,
                             town VARCHAR(100) NOT NULL,
                             country_name VARCHAR(100) NOT NULL,
                             time_zone VARCHAR(50) NOT NULL,
                             is_hq BOOLEAN NOT NULL
);
