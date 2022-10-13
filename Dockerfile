FROM golang:alpine AS builder

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Move working directory to /build
WORKDIR /build

# Copy the code into the container
COPY . .
RUN go mod download

WORKDIR /build/cmd
# Build the application
RUN go build -o ha 

# Build a small image
FROM scratch

COPY --from=builder /build/cmd/ha /

# Command to run
ENTRYPOINT ["/ha"]
CMD [ "-port=:8080", "-sid=" ]
