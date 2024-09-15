FROM golang:latest

WORKDIR /usr/src/storage

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o app ./services/storage/cmd/