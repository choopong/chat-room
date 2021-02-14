#!/bin/sh

source .env

heroku container:login
heroku container:push web --app=$HEROKU_APP
heroku config:set GIN_MODE=release --app=$HEROKU_APP
heroku container:release web --app=$HEROKU_APP