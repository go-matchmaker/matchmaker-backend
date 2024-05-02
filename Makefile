HTTP_MAIN_PACKAGE_PATH := cmd/http

PROTO_DIR := proto
PB_DIR := proto/pb
PROTOC := protoc
GRPC_PLUGIN := protoc-gen-go
GRPC_GATEWAY_PLUGIN := protoc-gen-grpc-gateway
PROTOC_OPTS := -I$(PROTO_DIR) --go_out=$(PB_DIR) --go_opt=paths=source_relative --go-grpc_out=$(PB_DIR) --go-grpc_opt=paths=source_relative

MIGRATION_FOLDER := internal/adapter/storage/postgres/migration
DB_URL := postgresql://postgres:P@ssw0rd@localhost:5432/postgres?sslmode=disable
# ==================================================================================== #
# DB
# ==================================================================================== #
create-db:
	docker exec -it postgres createdb --username=root --owner=root simple_bank

drop-db:
	docker exec -it postgres dropdb postgres

migrat-eup:
	migrate -path "$(MIGRATION_FOLDER)" -database "$(DB_URL)" -verbose up

migrate-up1:
	migrate -path "$(MIGRATION_FOLDER)" -database "$(DB_URL)" -verbose up 1

migrate-down:
	migrate -path "$(MIGRATION_FOLDER)" -database "$(DB_URL)" -verbose down

migrate-down1:
	migrate -path "$(MIGRATION_FOLDER)" -database "$(DB_URL)" -verbose down 1

new-migration:
	migrate create -ext sql -dir "$(MIGRATION_FOLDER)" -seq $(name)

migrate-force:
	migrate -path $(MIGRATION_FOLDER) -database "$(DB_URL)" force 1

db-docs:
	dbdocs build doc/db.dbml

# ==================================================================================== #
# HELPERS
# ==================================================================================== #
## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

.PHONY: no-dirty
no-dirty:
	git diff --exit-code

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #
## format: format code
.PHONY: format
format:
	find . -name '*.go' -exec gofumpt -w {} +

## tidy: format code and tidy modfile
.PHONY: tidy
tidy:
	go fmt ./...
	go mod tidy -v

## audit: run quality control checks
.PHONY: audit
audit:
	go mod verify
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...
	go test -race -buildvcs -vet=off ./...

## linter: run linters
.PHONY: linter-golangci
linter-golangci: ### check by golangci linter
	golangci-lint run

## clean: clean-up
.PHONY: clean
clean:
	go clean

#sec: sec
.PHONY: sec
sec:
	gosec ./...

#critic: critic
.PHONY: critic
critic:
	gocritic check -enableAll ./...

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #
## test: run all tests
.PHONY: test
test:
	go test -v -race -buildvcs ./...

## test/cover: run all tests and display coverage
.PHONY: test/cover
test/cover:
	go test -v -race -buildvcs -coverprofile=/tmp/coverage.out ./...
	go tool cover -html=/tmp/coverage.out

## build: build the application
.PHONY: build
build:
	# Include additional build steps, like TypeScript, SCSS, or Tailwind compilation here...
	go build -o=/tmp/bin/${BINARY_NAME} ${MAIN_PACKAGE_PATH}

## run: run the application
.PHONY: run
run: build
	/tmp/bin/${BINARY_NAME}

## run/live: run the application with reloading on file changes
.PHONY: run/live
run/live:
	go run github.com/cosmtrek/air@v1.43.0 \
		--build.cmd "make build" --build.bin "/tmp/bin/${BINARY_NAME}" --build.delay "100" \
		--build.exclude_dir "" \
		--build.include_ext "go,tpl,tmpl,html,css,scss,js,ts,sql,jpeg,jpg,gif,png,bmp,svg,webp,ico" \
		--misc.clean_on_exit "true"

# ==================================================================================== #
# RUN
# ==================================================================================== #
## run: run the applications
run: run-http

.PHONY: run-http
run-http:
	cd $(HTTP_MAIN_PACKAGE_PATH) && go mod tidy && go mod download && \
    CGO_ENABLED=0 go run github.com/go-matchmaker/matchmaker-server/$(HTTP_MAIN_PACKAGE_PATH)

# ==================================================================================== #
# Docker
# ==================================================================================== #
## docker-compose: run docker-compose
docker-compose: docker-compose-stop docker-compose-start
.PHONY: docker-compose

.PHONY: docker-compose-start
docker-compose-start:
	docker-compose up --build

.PHONY: docker-dependency-start
docker-dependency-start:
	docker-compose -f docker-compose-core.yaml up -d

.PHONY: docker-compose-stop
docker-compose-stop:
	ddocker-compose down

# ==================================================================================== #
# Generated Files
# ==================================================================================== #
## sqlc-generate: generate sqlc
.PHONY: sqlc-generate
sqlc-generate:
	sqlc generate

.PHONY: sqlc-clean
sqlc-clean:
	rm -rf internal/adapter/storage/postgres/sqlc/generated

## wire: generate wire
.PHONY: wire-generate
wire-generate:
	cd internal/adapter/app && wire && cd -

.PHONY: wire-clean
wire-clean:
	cd internal/adapter/app && rm wire_gen.go && cd -

## proto: generate protobuf files
.PHONY: proto-generate
proto-generate:
	$(PROTOC) $(PROTOC_OPTS) $(PROTO_DIR)/*.proto

.PHONY: proto-clean
proto-clean:
	rm proto/pb/*.pb.go;

.PHONY: dependency-generate
dependency-generate: wire-generate sqlc-generate

.PHONY: dependency-clean
dependency-clean: wire-clean sqlc-clean