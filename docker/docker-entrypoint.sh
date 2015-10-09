#!/bin/bash

# in the vein of
# https://github.com/docker-library/logstash/blob/master/1.5/docker-entrypoint.sh

set -e

if [ ! -f /etc/SCAMP/soa.conf ]; then
	init-system-config
fi

if [ ! -f /etc/SCAMP/services/logging.key ]; then
	provision-soa-service logging main
fi

ELASTICSEARCH_HOST="localhost"
ELASTICSEARCH_PORT=9200

until nc -z $ELASTICSEARCH_HOST $ELASTICSEARCH_PORT; do
    echo "$(date) - waiting for elasticsearch at $ELASTICSEARCH_HOST:$ELASTICSEARCH_PORT..."
    sleep 1
done

exec /usr/bin/hekad "-config=/etc/heka/conf.d"
