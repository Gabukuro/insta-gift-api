default: run

init:
	go mod init github.com/Gabukuro/insta-gift-api && \
	go mod tidy && \
	go mod vendor

# Install dependencies
install:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.54.2 && \
	go install github.com/golang/mock/mockgen@latest && \
	go install gotest.tools/gotestsum@latest && \
	go mod tidy && \
	go mod vendor

# Upgrade packages to last version
upgrade-pkgs:
	go get -u ./... && make install

# Build the cover HTML file
cover:
	go tool cover -html=coverage.out -o coverage.html

# Build
build: 
	go build cmd/main.go

# Build production
build-prod:
	go build -o /tmp/main cmd/main.go

# Run
run:
	cd scripts/db/ && ./setup_db.sh && cd ../../ && LOGGER_FORMAT=cli BUNDEBUG=2 go run cmd/main.go

# Run dev
run-dev:
	nodemon --exec go run cmd/main.go --signal SIGTERM

# Run lint
lint:
	golangci-lint run ./...

# Run tests
test: test-all

# Run unit tests
test-unit:
	gotestsum --format testname -- $(shell go list ./... | egrep '/internal' | egrep -v '/mock|/testhelper') -coverprofile=coverage.out -coverpkg=./internal/... && \
	egrep -v 'mock|testhelper' coverage.out > tmpcoverage && mv tmpcoverage coverage.out

# Run e2e tests
test-e2e:
	gotestsum --format testname -- ./test/integration/... -coverprofile=coverage.out -coverpkg=./internal/... && \
	egrep -v 'mock|testhelper' coverage.out > tmpcoverage && mv tmpcoverage coverage.out

# Run all tests
test-all:
	(gotestsum --format testname -- ./... -coverprofile=coverage.out -coverpkg=./internal/... && \
	egrep -v 'mock|testhelper' coverage.out > tmpcoverage && mv tmpcoverage coverage.out) && \
	gotestsum --format testname -- -json ./... | go-junit-report > report-all.xml

