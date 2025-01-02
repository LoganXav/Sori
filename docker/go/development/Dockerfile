FROM golang:1.19.5-alpine3.17 AS build

# Install necessary build tools including gcc
RUN apk update && apk add --no-cache \
    make \
    git \
    gcc \
    musl-dev \
    ca-certificates \
    wget && \
    update-ca-certificates

# Create a non-root user
ENV USER=appuser
ENV UID=10001
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"

WORKDIR $GOPATH/src/sori

COPY . .

# Download Go dependencies
RUN go mod download

# Install additional tools (swag, reflex, golangci-lint)
RUN wget https://github.com/swaggo/swag/releases/download/v1.8.10/swag_1.8.10_Linux_x86_64.tar.gz && \
    tar -xvzf swag_1.8.10_Linux_x86_64.tar.gz && \
    mv swag /go/bin

RUN wget https://github.com/cespare/reflex/releases/download/v0.3.1/reflex_linux_amd64.tar.gz && \
    tar -xvzf reflex_linux_amd64.tar.gz && \
    mv reflex_linux_amd64/reflex /go/bin


RUN wget https://github.com/golangci/golangci-lint/releases/download/v1.51.1/golangci-lint-1.51.1-linux-amd64.tar.gz && \
    tar -xvzf golangci-lint-1.51.1-linux-amd64.tar.gz && \
    mv golangci-lint-1.51.1-linux-amd64/golangci-lint /go/bin

# Build the Go application
RUN go build -o server -gcflags=all=-l -ldflags="-w -s" main.go

# Expose application port
EXPOSE 3000

# Run the application in watch mode
CMD ["make", "watch"]
