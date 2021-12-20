## Build
FROM golang:1.17.2-alpine AS builder

ARG APP_NAME=dms-be

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 go build -o ./$APP_NAME

# Run
FROM alpine:3.14
WORKDIR /app
COPY --from=builder /app/$APP_NAME .

EXPOSE 8080

CMD /app/dms-be