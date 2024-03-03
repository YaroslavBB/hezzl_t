#!/bin/bash
export ROOT=..
export CONF_PATH=$ROOT/config/config.yaml
export MIGRATION_LOCK=$ROOT/migration/migrations.lock
export INIT_PG_DB=$ROOT/scripts/init_postgres_db.sh
export DOCKER_PATH=../../docker/ 