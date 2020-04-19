#!/usr/bin/env sh

docker run -d -p 127.0.0.1:3306:3306 --name shoppingcart_mysql \
    -e MYSQL_ROOT_PASSWORD=password \
    -e MYSQL_DATABASE=shoppingcart \
    -e MYSQL_USER=shoppingcart \
    -e MYSQL_PASSWORD=secret \
    -e MYSQL_TCP_PORT=3306 \
    mysql:5.7 \
    --innodb_log_file_size=256MB \
    --innodb_buffer_pool_size=512MB \
    --max_allowed_packet=16MB \
    --local-infile=1