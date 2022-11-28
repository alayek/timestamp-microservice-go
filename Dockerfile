
FROM golang:1.19.3-alpine as Builder

ARG GIT_VERSION

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY ./build.sh ./

COPY *.go ./

RUN GIT_VERSION=${GIT_VERSION} sh ./build.sh

FROM alpine as Final

COPY --from=Builder /app/timestamp /sbin/timestamp

EXPOSE 8080

ENTRYPOINT [ "/sbin/timestamp" ]