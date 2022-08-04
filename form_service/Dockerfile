FROM golang:alpine

ENV GIN_MODE=release
ENV PORT=8080

WORKDIR /go/src/formService

COPY . /go/src/formService

# Run the two commands below to install git and dependencies for the project.
RUN apk update && apk add --no-cache git
RUN go mod tidy

EXPOSE $PORT

CMD go run server.go

# ENTRYPOINT ["./app"]