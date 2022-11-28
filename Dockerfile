FROM golang:1.19.3-alpine as Builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY *.go ./

RUN go build -o /timestamp

FROM alpine:edge as Final

COPY --from=Builder /timestamp /sbin/timestamp

EXPOSE 8080

ENTRYPOINT [ "/sbin/timestamp" ]