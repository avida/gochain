#!/bin/bash
PORT=5432
POSTRESQL_PASS=1888
USER=$1
PASSWORD=$2
usage () {
  echo "Usage is"
}
if [ "$USER" == "" ] || [ "$PASSWORD" == "" ]
then
  usage
  exit 0
fi
echo "user is $USER"
echo "password is $PASSWORD"
export PGPASSWORD=$POSTRESQL_PASS 

# Add user
echo "CREATE USER $USER WITH PASSWORD '$PASSWORD';" | \
psql  -p $PORT -h localhost -U postgres  -f - 

# Create db
echo "CREATE DATABASE $USER wITH OWNER=$USER;" | \
psql  -p $PORT -h localhost -U postgres  -f - 

# create tables
PGPASSWORD=$PASSWORD psql  -p $PORT -h localhost -U $USER -d $USER -f db/schema.sql

