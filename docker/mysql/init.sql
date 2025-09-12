-- Membuat database pertama
CREATE DATABASE IF NOT EXISTS first_go_api_db;

-- Membuat database kedua
CREATE DATABASE IF NOT EXISTS first_go_api_db_test;

-- (Opsional tapi sangat direkomendasikan)
-- Membuat user baru dan memberikan hak akses ke kedua database tersebut
-- CREATE USER 'encang-cutbray'@'%' IDENTIFIED BY 'secret214';
GRANT ALL PRIVILEGES ON `first_go_api_db`.* TO 'encang-cutbray'@'%';
GRANT ALL PRIVILEGES ON `first_go_api_db_test`.* TO 'encang-cutbray'@'%';
FLUSH PRIVILEGES;