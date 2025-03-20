FROM golang:1.24 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download
RUN go install github.com/openshift-online/ocm-cli/cmd/ocm@latest && mv /go/bin/ocm /app/

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -o osplc

FROM registry.access.redhat.com/ubi9/ubi-micro

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/osplc /app/osplc
COPY --from=builder /app/ocm /app/ocm