#!/bin/bash

name='baur-local-db'

if [ $1 == "up" ]; then
  if ! docker ps --format '{{.Names}}' | grep -w $name &> /dev/null; then
      if [ "$(docker ps -aq -f status=exited -f name=$name)" ]; then
          echo "Container $name already exists, starting.."
          docker start $name
          exit 0
      fi

      echo "Initializing container $name"
      docker run -d --name $name -p 15432:5432 -e POSTGRES_PASSWORD=baur -e POSTGRES_DB=baur postgres:12
      docker run -e SLEEP_LENGTH=1 -e TIMEOUT_LENGTH=60 dadarek/wait-for-dependencies host.docker.internal:15432
      sleep 5
      baur init db
      echo "Done"
  fi
fi

if [ $1 == "down" ]; then
  docker rm -f $name || 0
fi