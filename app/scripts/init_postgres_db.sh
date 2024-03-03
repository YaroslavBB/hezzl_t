#!/bin/bash

source variables.sh

sql_file=../app/migration/init_postgres_db.sql

cd $DOCKER_PATH

docker compose exec -T db psql -h localhost -p 5432 -U postgres -W < $sql_file
