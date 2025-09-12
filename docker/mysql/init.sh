#!/bin/bash
set -e

mysql -v -uroot -p"${MYSQL_ROOT_PASSWORD}" <<-EOSQL
    -- Membuat database utama dari environment variable
    CREATE DATABASE IF NOT EXISTS \`${DB_DATABASE}\`;

    -- Membuat database tes dari environment variable
    CREATE DATABASE IF NOT EXISTS \`${DB_DATABASE_TESTING}\`;

    -- Memberikan hak akses kepada user yang dibuat oleh Docker
    GRANT ALL PRIVILEGES ON \`${DB_DATABASE}\`.* TO '${MYSQL_USER}'@'%';
    GRANT ALL PRIVILEGES ON \`${DB_DATABASE_TESTING}\`.* TO '${MYSQL_USER}'@'%';
    FLUSH PRIVILEGES;
EOSQL