#!/bin/sh

if [ ! -f "/config/oauth.json" ]; then
    cp /src/config/oauth_example.json /config/oauth.json
fi

exec "$@"