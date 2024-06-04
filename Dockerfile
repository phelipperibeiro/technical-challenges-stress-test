# Use the official Golang image
FROM golang:1.22 as builder

# Instalar os certificados raiz
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

# Atualizar os certificados
RUN update-ca-certificates

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux COARCH=amd64 go build -o main main.go

# Use a minimal Docker image to run the Go app
FROM scratch

# Copiar os certificados raiz para a nova imagem
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .

# Set the entry point to the executable
ENTRYPOINT ["./main"]


# docker build --no-cache  -t technical-challenges-stress-tes .
# docker run technical-challenges-stress-tes --url=http://globo.com --requests=30 --concurrency=20
# docker run technical-challenges-stress-tes --url=http://fullcycle.com.br --requests=1000 --concurrency=100