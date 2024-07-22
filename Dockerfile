FROM golang:1.21

WORKDIR /app

# Download Go modules
COPY . .
RUN go mod download

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o reporter
ENV GIN_MODE=release

EXPOSE 8080

# Run
ENTRYPOINT [ "./reporter", "-syslog=false" ]