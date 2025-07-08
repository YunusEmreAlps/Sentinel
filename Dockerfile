# ----------------------------
# STAGE-1: build stage
FROM golang:1.21-alpine3.17 AS build-env
RUN apk add build-base

WORKDIR /src

COPY . .

# set build option to linux/amd64 for decrease build size
RUN CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    go build -o main .

# ----------------------------
# STAGE-2: output stage
FROM alpine:3.17

WORKDIR /app

# copy necessary files from build stage
COPY --from=build-env /src/main /app/

# create user and give it permissions
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
RUN chown -R appuser:appgroup /app
USER appuser

ENTRYPOINT ./main
