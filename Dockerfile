# Build container
FROM --platform=$BUILDPLATFORM golang:1.22.5-alpine AS build

ARG VERSION
ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

# Install build dependencies (e.g. gcc)
RUN apk add build-base

COPY ./src ./src
COPY ./go.mod .
COPY ./go.sum .

# Download Go dependencies
RUN go mod download
# Generate Templ templates
RUN go install github.com/a-h/templ/cmd/templ@latest
RUN templ generate

# Compile with CGo enabled
RUN CGO_ENABLED=1 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o johtotimes ./src

# Deployment container
FROM alpine:latest

WORKDIR /app

COPY --from=build /app/johtotimes .

EXPOSE 3000

ENTRYPOINT ["/app/johtotimes"]
