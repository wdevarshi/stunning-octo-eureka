.PHONY: build build-alpine clean test default install generate run run-docker stress-test stress-test-10qps stress-test-100qps stress-test-check populate-data pdf-report

BIN_NAME=myapp

SHELL := $(shell which bash)
VERSION := $(shell grep "const Version " version/version.go | sed -E 's/.*"(.+)"$$/\1/')
GIT_COMMIT=$(shell git rev-parse HEAD)
GIT_DIRTY=$(shell test -n "`git status --porcelain`" && echo "+DIRTY" || true)
BUILD_DATE=$(shell date '+%Y-%m-%d-%H:%M:%S')
IMAGE_NAME := "myapp"

default: test

help:
	@echo 'Management commands for myapp:'
	@echo
	@echo 'Usage:'
	@echo '    make bench               Run benchmarks.'
	@echo '    make build               Compile the project and generate a binary.'
	@echo '    make build-docker        Build a docker image.'
	@echo '    make clean               Clean the directory tree.'
	@echo '    make coverage-html       Generate test coverage report.'
	@echo '    make dep                 Update dependencies.'
	@echo '    make generate            Generate code from proto files.'
	@echo '    make help                Show this message.'
	@echo '    make lint                Run linters on the project.'
	@echo '    make mock                Generate mocks for interfaces.'
	@echo '    make pdf-report          Generate PDF report from LaTeX.'
	@echo '    make populate-data       Populate database with sample data.'
	@echo '    make run                 Run the project locally.'
	@echo '    make run-docker          Run the project in a docker container.'
	@echo '    make runj                Run the project locally with jq log parsing.'
	@echo '    make stress-test         Run all stress tests (10 QPS + 100 QPS).'
	@echo '    make stress-test-10qps   Run 10 QPS stress test.'
	@echo '    make stress-test-100qps  Run 100 QPS stress test.'
	@echo '    make test                Run tests.'
	@echo

build:
	@echo "building ${BIN_NAME} ${VERSION}"
	@echo "GOPATH=${GOPATH}"
	go build -ldflags "-X github.com/bluesg/transport-analytics/version.GitCommit=${GIT_COMMIT}${GIT_DIRTY} -X github.com/bluesg/transport-analytics/version.BuildDate=${BUILD_DATE} -X github.com/bluesg/transport-analytics/version.Branch=${GIT_BRANCH}" -o bin/${BIN_NAME}

build-alpine:
	@echo "building ${BIN_NAME} ${VERSION}"
	@echo "GOPATH=${GOPATH}"
	go build -ldflags '-w -linkmode external -extldflags "-static" -X github.com/bluesg/transport-analytics/version.GitCommit=${GIT_COMMIT}${GIT_DIRTY} -X github.com/bluesg/transport-analytics/version.BuildDate=${BUILD_DATE} -X github.com/bluesg/transport-analytics/version.Branch=${GIT_BRANCH} ' -o bin/${BIN_NAME}

build-docker:
	@echo "building image ${BIN_NAME} ${VERSION} $(GIT_COMMIT)"
	docker build --build-arg VERSION=${VERSION} --build-arg GIT_COMMIT=$(GIT_COMMIT) --build-arg GIT_BRANCH=$(GIT_BRANCH) -t $(IMAGE_NAME):local .

dep:
	go mod tidy

install:
	go install \
		github.com/vektra/mockery/v2 \
		github.com/bufbuild/buf/cmd/buf \
		github.com/golangci/golangci-lint/cmd/golangci-lint

generate: install
	buf generate --path proto/*.proto

clean:
	@test ! -e bin/${BIN_NAME} || rm bin/${BIN_NAME}
	go clean ./...

test:
	go test -race -coverpkg=.,./config/...,./service/... -coverprofile cover.out ./...
	go tool cover -func=cover.out

coverage-html:
	go tool cover -html=cover.out -o=cover.html

bench:
	# -run=^B negates all tests
	go test -bench=. -run=^B -benchtime 10s -benchmem ./...

lint: install
	golangci-lint run --timeout 5m

mock: install
	mockery --config .mockery.yaml

run: build
	@echo
	@echo "swagger ui available at http://localhost:9091/swagger/"
	@echo
	@set -a; source local.env; ./bin/${BIN_NAME}

runj: build
	@echo
	@echo "swagger ui available at http://localhost:9091/swagger/"
	@echo
	@set -a; source local.env && ./bin/${BIN_NAME} 1> >(jq -R "fromjson? | ." -C)

run-docker: build-docker
	docker run -p 9091:9091 -p 9090:9090 --env-file local.env ${IMAGE_NAME}:local

stress-test-check:
	@command -v vegeta >/dev/null 2>&1 || { echo "Error: vegeta is not installed. Install with: brew install vegeta (macOS) or go install github.com/tsenart/vegeta@latest"; exit 1; }
	@curl -s http://localhost:9091/health > /dev/null 2>&1 || { echo "Error: API is not running at http://localhost:9091. Please start the application first with: make run or docker-compose up"; exit 1; }

stress-test-10qps: stress-test-check
	@echo "Running 10 QPS stress test..."
	@cd scripts/stress-test && ./run-10qps.sh

stress-test-100qps: stress-test-check
	@echo "Running 100 QPS stress test..."
	@cd scripts/stress-test && ./run-100qps.sh

stress-test: stress-test-check
	@echo "Running all stress tests..."
	@cd scripts/stress-test && ./run-all-tests.sh

populate-data:
	@curl -s http://localhost:9091/health > /dev/null 2>&1 || { echo "Error: API is not running at http://localhost:9091. Please start the application first with: make run or docker-compose up"; exit 1; }
	@echo "Populating database with sample data..."
	@cd scripts/populate-data && go run main.go

pdf-report:
	@echo "Generating PDF report from LaTeX..."
	@if command -v pdflatex >/dev/null 2>&1; then \
		echo "Using local pdflatex..."; \
		pdflatex -interaction=nonstopmode PROJECT_REPORT.tex > /dev/null 2>&1; \
		pdflatex -interaction=nonstopmode PROJECT_REPORT.tex > /dev/null 2>&1; \
		rm -f PROJECT_REPORT.aux PROJECT_REPORT.log PROJECT_REPORT.out PROJECT_REPORT.toc; \
		echo "✓ PDF generated: PROJECT_REPORT.pdf"; \
	else \
		echo "pdflatex not found, using Docker..."; \
		docker run --rm -v "$(PWD):/workspace" -w /workspace texlive/texlive:latest pdflatex -interaction=nonstopmode PROJECT_REPORT.tex > /dev/null 2>&1; \
		docker run --rm -v "$(PWD):/workspace" -w /workspace texlive/texlive:latest pdflatex -interaction=nonstopmode PROJECT_REPORT.tex > /dev/null 2>&1; \
		rm -f PROJECT_REPORT.aux PROJECT_REPORT.log PROJECT_REPORT.out PROJECT_REPORT.toc; \
		echo "✓ PDF generated: PROJECT_REPORT.pdf"; \
	fi
