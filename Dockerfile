
#
## Dependencies Stage
FROM golang:1.19.3-alpine3.15 AS dependencies-stage
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

#
## Build Stage
FROM dependencies-stage AS build-stage
COPY . .
RUN go build -o cinnox-homework 

#
## Release Stage
FROM alpine:3.15 AS release-stage
ENV GIN_MODE=release

COPY --from=build-stage /app/.env /.env
COPY --from=build-stage /app/cinnox-homework /cinnox-homework

EXPOSE 8080

CMD [ "" ]
ENTRYPOINT [ "/cinnox-homework" ]