BINARY_NAME=reporter
MIGRATE_BINARY_NAME=migrate

all: lint test build-linux build-deb build-rpm

build:
	go build -o ${BINARY_NAME} cmd/reporter/main.go

build-linux:
	CGO_ENABLED=0 \
	GOARCH=amd64 \
	GOOS=linux \
	go build -o ${BINARY_NAME} cmd/reporter/main.go

build-migrate:
	go build -o ${MIGRATE_BINARY_NAME} cmd/migrate/main.go

run:
	go build -o ${BINARY_NAME} cmd/reporter/main.go
	./${BINARY_NAME}

dev:
	go run cmd/reporter/main.go

test:
	go test -v ./...

lint:
	golangci-lint run -c .golanci.yml 

build-deb:
	cd scripts && ./create-deb-package.sh

build-rpm:
	cd scripts && ./create-rpm-package.sh
 
clean:
	go clean
	rm ${BINARY_NAME}