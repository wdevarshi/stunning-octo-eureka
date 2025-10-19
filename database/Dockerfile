FROM postgres:15-alpine

COPY init.sql /docker-entrypoint-initdb.d/

ENV POSTGRES_DB=transport_reliability
ENV POSTGRES_USER=lta_user
ENV POSTGRES_PASSWORD=lta_pass_2025
