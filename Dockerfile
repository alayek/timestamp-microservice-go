
FROM golang:1.19.3-alpine as Builder

RUN apk update && apk add git

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY ./build.sh ./

COPY ./.git ./

COPY *.go ./

RUN sh ./build.sh

FROM alpine as Final

COPY --from=Builder /app/timestamp /sbin/timestamp

EXPOSE 8080

ENTRYPOINT [ "/sbin/timestamp" ]