# ==================================================================================== #
# BASE IMAGE
# ==================================================================================== #
FROM postgres:16.3-alpine as base

COPY --chown=postgres:postgres postgres.conf /etc/postgresql/postgres.conf
COPY --chown=postgres:postgres create-databases.sh /docker-entrypoint-initdb.d

EXPOSE 5432

CMD ["-c", "config_file=/etc/postgresql/postgres.conf"]
