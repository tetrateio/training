#!/bin/sh

set -e

# we don't want to replace undefined env vars.
envsubst "$(env | sed -e 's/=.*//' -e 's/^/\$/g')" \
    < /etc/nginx/conf.d/nginx.conf.template > /etc/nginx/conf.d/default.conf

exec nginx -g 'daemon off;'
