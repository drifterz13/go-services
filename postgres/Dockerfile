FROM postgres:13.4-alpine

COPY init.sql /docker-entrypoint-initdb.d/init.sql

CMD ["postgres"]