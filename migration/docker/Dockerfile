FROM alpine:3.22.0

RUN wget https://github.com/amacneil/dbmate/releases/latest/download/dbmate-linux-amd64 -O /usr/bin/dbmate
RUN chmod +x /usr/bin/dbmate
RUN apk add mc

COPY ./script /script
COPY ./migration /migration

RUN chmod +x /script/startMigration.sh


WORKDIR /script

ENTRYPOINT ["/script/startMigration.sh"]
