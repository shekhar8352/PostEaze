FROM golang:1.23-alpine

# Install git (sometimes needed for go packages)
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go.mod and vendor directory
COPY go.mod go.sum ./
COPY vendor/ ./vendor/

# Copy the rest of the source code
COPY . .

# Build with vendor mode
RUN go build -mod=vendor -o main .

EXPOSE 8000

CMD ["./main"]
