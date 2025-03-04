CREATE TABLE swiftTable (
                             country_iso2 CHAR(2) NOT NULL,
                             code VARCHAR(11) NOT NULL UNIQUE,
                             bank_name VARCHAR(255) NOT NULL,
                             address VARCHAR(255) NOT NULL,
                             country_name VARCHAR(100) NOT NULL,
                             is_hq BOOLEAN NOT NULL
);
 