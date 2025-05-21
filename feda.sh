#!/bin/bash

docker exec $1 /bin/sh -c "kill -2 27"
sleep 10
docker cp $1:/usr/local/bin/flamegraph.svg $2.svg
