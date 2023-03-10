FROM golang:1.20 AS build
WORKDIR /go/src

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

ENV CGO_ENABLED=0
RUN go get -d -v ./...

RUN go build -a -installsuffix cgo -o discord *.go

# Restarting from scratch loses updated ca certificates. We need these for the discord client 
FROM buildpack-deps:bullseye-scm
RUN apt update -y
RUN update-ca-certificates

COPY --from=build /go/src/discord ./
ENTRYPOINT ["./discord"]
