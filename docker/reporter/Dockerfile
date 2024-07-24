FROM golang:1.21

WORKDIR /app

# Download Go modules
COPY . .
RUN go mod download

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/reporter/reporter cmd/reporter/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/migrate/migrate cmd/migrate/main.go
ENV GIN_MODE=release

EXPOSE 8080

# Run
ENTRYPOINT [ "./bin/reporter/reporter", "-syslog=false" ]