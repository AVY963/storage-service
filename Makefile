# --- CONFIG ---

LOCAL_BIN     := $(CURDIR)/bin

MIGRATE_VERSION       := v4.18.2

define install_tool
	GOBIN=$(LOCAL_BIN) go install $(1)@$(2)
endef

# --- INSTALL TOOLS ---

.PHONY: install
install:
	mkdir -p $(LOCAL_BIN)
	# go mod tidy
	$(call install_tool,github.com/golang-migrate/migrate/v4/cmd/migrate,$(MIGRATE_VERSION))
	GOBIN=$(LOCAL_BIN) go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@$(MIGRATE_VERSION)

# --- DATABASE ---

.PHONY: db-up
db-up:
	docker-compose up -d postgres
.PHONY: db-down
db-down:
	docker-compose down

# --- MIGRATIONS ---

.PHONY: migrate-create
migrate-create:
	$(LOCAL_BIN)/migrate create -ext sql -dir migrations -seq $(name)

.PHONY: migrate-up
migrate-up:
	$(LOCAL_BIN)/migrate -path migrations -database "$(DB_URL)" up

.PHONY: migrate-down
migrate-down:
	$(LOCAL_BIN)/migrate -path migrations -database "$(DB_URL)" down

.PHONY: migrate-force
migrate-force:
	$(LOCAL_BIN)/migrate -path migrations -database "$(DB_URL)" force $(version)

# --- RUN ---

.PHONY: run
run:
	go run cmd/main.go

# --- BUILD ---

.PHONY: build
build:
	go build -o bin/fileserver cmd/main.go
