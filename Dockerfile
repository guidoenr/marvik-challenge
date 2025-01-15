# MULTI-STAGE Dockerfile
# STAGE 1: build the go binary
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY . .

# build the Go application (binary)
RUN go mod download
RUN go build -o app .

# STAGE 2: init postgresql and load it with data
FROM postgres:17.2-alpine3.21 AS base

# (this makes the image start with the init.sql script)
COPY init.sql /docker-entrypoint-initdb.d/

# default PSQL creds
ENV POSTGRES_USER=postgres
ENV POSTGRES_PASSWORD=mysecretpassword
ENV POSTGRES_DB=postgres

# DEFAULT DEBUG MODE (disable later)
ENV DEBUG=true

COPY --from=builder /app/app /usr/local/bin/app

EXPOSE 8080

# start the default entry point until pg_isready 
ENTRYPOINT ["sh", "-c", "docker-entrypoint.sh postgres & while ! pg_isready -h localhost -p 5432 -U postgres; do echo 'waiting for database...'; sleep 1; done; app"]


