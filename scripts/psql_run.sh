#!/bin/bash
usage () {
  echo "usage is psql_run.sh <container name> <postrgres password>"
}
if [ $# -ne 2 ]
  then
    usage 
    exit
  fi
CONTAINER=$1
PASSWORD=$2
EXISTS=$(docker ps -a | grep "\s$CONTAINER$")
IS_RUNNING=$(docker ps | grep "\s$CONTAINER$")
if [ "$EXISTS" == "" ]
  then
    echo "container $CONTAINER doesnt exist"
    docker run  --name $CONTAINER -e POSTGRES_PASSWORD=$PASSWORD -d -p 5432:5432 postgres
  fi

if [ "$IS_RUNNING" != "" ]
  then
    echo "$CONTAINER is already running"
  else
    echo "container $CONTAINER is not running"
    echo "Starting $CONTAINER"
    docker start $CONTAINER
  fi

