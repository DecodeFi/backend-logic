#!/bin/bash

# PGPASSWORD!
HOST="127.0.0.1"
DB=$1
USER=$2

psql --host $HOST --user $USER  --db $DB -f ./migration.sql
