FROM golang:1.19.3-alpine as Builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY ./build.sh ./

COPY *.go ./

RUN sh ./build.sh

FROM alpine:edge as Final

COPY --from=Builder /timestamp /sbin/timestamp

EXPOSE 8080

ENTRYPOINT [ "/sbin/timestamp" ]