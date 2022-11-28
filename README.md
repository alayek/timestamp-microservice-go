A Golang implementation of [FreeCodeCamp API project Timestamp Microservice](https://www.freecodecamp.org/learn/back-end-development-and-apis/back-end-development-and-apis-projects/timestamp-microservice)

## How to Run Locally

### With Docker

- Install Docker
- Run `docker build -t timestamp .` to build the docker image.
- Now run the docker image with `docker run -p 8080:8080 --rm -d --name ts timestamp`
- Once the docker image has a running instance, make API calls to localhost as specified in the problem statement

### Without Docker

- Install Go and set up `GOROOT`, `GOPATH` properly.
- Run `go mod download`
- Run Go build with `go build main.go -o timestamp`
- Execute the generated binary with `./timestamp`
