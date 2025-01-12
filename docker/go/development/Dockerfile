# Use the official Go image based on Alpine for minimal footprint
FROM golang:1.23-alpine AS build

# Install necessary build tools, including GCC and OpenJDK 11 for FastQC
RUN apk update && apk add --no-cache \
    make \
    git \
    gcc \
    musl-dev \
    ca-certificates \
    wget \
    perl \
    bash \
    openjdk11 && \
    update-ca-certificates

# Install FastQC and ensure it's executable
RUN wget -O /tmp/fastqc.zip https://www.bioinformatics.babraham.ac.uk/projects/fastqc/fastqc_v0.11.9.zip && \
    unzip /tmp/fastqc.zip -d /opt && \
    chmod +x /opt/FastQC/fastqc && \
    ln -s /opt/FastQC/fastqc /go/bin

# Set JAVA_HOME for FastQC's runtime environment and update the PATH
ENV JAVA_HOME=/usr/lib/jvm/java-11-openjdk
ENV PATH="/go/bin:${PATH}"

# Create a non-root user to enhance security with a writable home directory
ENV USER=appuser
ENV UID=10001
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/home/${USER}" \
    --shell "/sbin/nologin" \
    --uid "${UID}" \
    "${USER}"

# Ensure the non-root user has ownership of relevant directories
RUN chown -R ${USER}:${USER} /opt/FastQC /go /home/${USER}

# Switch to the non-root user
USER ${USER}

# Set Go build cache to a writable directory
ENV GOCACHE=/home/${USER}/.cache/go-build

# Set the working directory for the Go application
WORKDIR $GOPATH/src/sori

# Verify FastQC installation
RUN fastqc --version

# Copy the application source code into the container
COPY . .

# Download Go dependencies to cache them during build
RUN go mod download

# Install additional tools: swag for generating Swagger docs, reflex for live reload, and golangci-lint for linting
RUN wget https://github.com/swaggo/swag/releases/download/v1.8.10/swag_1.8.10_Linux_x86_64.tar.gz && \
    tar -xvzf swag_1.8.10_Linux_x86_64.tar.gz && \
    mv swag /go/bin

RUN wget https://github.com/cespare/reflex/releases/download/v0.3.1/reflex_linux_amd64.tar.gz && \
    tar -xvzf reflex_linux_amd64.tar.gz && \
    mv reflex_linux_amd64/reflex /go/bin

RUN wget https://github.com/golangci/golangci-lint/releases/download/v1.51.1/golangci-lint-1.51.1-linux-amd64.tar.gz && \
    tar -xvzf golangci-lint-1.51.1-linux-amd64.tar.gz && \
    mv golangci-lint-1.51.1-linux-amd64/golangci-lint /go/bin

# Build the Go application binary
RUN go build -o server -gcflags=all=-l -ldflags="-w -s" main.go

# Expose the application's listening port
EXPOSE 3000

# Set the command to run the application in watch mode using Makefile
CMD ["make", "watch"]
