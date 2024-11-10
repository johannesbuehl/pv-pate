CREATE TABLE elements (mid CHAR(6) NOT NULL KEY , name TINYTEXT NOT NULL DEFAULT "", mail TINYTEXT, reservation TIMESTAMP NULL DEFAULT current_timestamp());
CREATE TABLE users (uid INT NOT NULL KEY auto_increment, name TINYTEXT NOT NULL, password binary(60) NOT NULL, tid INT NOT NULL DEFAULT 0);
