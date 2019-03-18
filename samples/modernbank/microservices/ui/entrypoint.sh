#!/bin/sh

set -e

if [[ -z "${API_HOST}" ]]; then
    echo "running without proxying..."
    cp /etc/nginx/conf.d/nginx.conf.template /etc/nginx/conf.d/default.conf
else
    # We don't want to replace undefined env vars.
    envsubst "$(env | sed -e 's/=.*//' -e 's/^/\$/g')" \
        < /etc/nginx/conf.d/nginx-with-proxy-pass.conf.template > /etc/nginx/conf.d/default.conf
fi

exec nginx -g 'daemon off;'
