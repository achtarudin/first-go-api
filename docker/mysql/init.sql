-- Membuat database utama dari environment variable
CREATE DATABASE IF NOT EXISTS `first_go_api_db`;

-- Membuat database tes dari environment variable
CREATE DATABASE IF NOT EXISTS `first_go_api_db_test`;

-- Memberikan hak akses kepada user yang dibuat oleh Docker
GRANT ALL PRIVILEGES ON `first_go_api_db`.* TO 'encang-cutbray'@'%';
GRANT ALL PRIVILEGES ON `first_go_api_db_test`.* TO 'encang-cutbray'@'%';
FLUSH PRIVILEGES;