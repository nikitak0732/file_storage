#!/bin/sh

set -e

echo Старт миграций

for item in $(find /migration/ -type d | sort);
do
    if  find ${item} -maxdepth 1 -name "*.sql" -print -quit | grep -q .; then
        echo ${item};
        dbmate --wait -u "postgres://$POSTGRES_USER:$POSTGRES_PASSWORD@$POSTGRES_HOST:$POSTGRES_PORT/$POSTGRES_DB?sslmode=disable" --migrations-dir "${item}" migrate
    fi
done

echo Миграции выполнены