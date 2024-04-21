# syntax=docker/dockerfile:1
FROM golang:1.22

################################################################################
# Create a stage for building the application.
ARG GO_VERSION=1.22.0
# FROM --platform=$BUILDPLATFORM golang:${GO_VERSION} AS build
WORKDIR /src/johtotimes

COPY go.mod go.sum ./

RUN go mod download
RUN go install github.com/cosmtrek/air@v1.51.0
RUN go install github.com/a-h/templ/cmd/templ@v0.2.590

# COPY . .
# RUN go build -o main ./src

EXPOSE 3000

CMD ["air"]
