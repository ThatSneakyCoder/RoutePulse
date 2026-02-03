# -------------------------
# Migration configuration
# -------------------------
MIGRATIONS_PATH := services/identity-service/internal/migrate/migrations

# Development-only DB address (via kubectl port-forward)
DB_ADDR := postgres://guest:guest@localhost:15432/routepulse?sslmode=disable

# -------------------------
# Migrations
# -------------------------
.PHONY: migrate-create
migrate-create:
	@migrate create -seq -ext sql -dir $(MIGRATIONS_PATH) $(filter-out $@,$(MAKECMDGOALS))

.PHONY: migrate-up
migrate-up:
	@migrate -path=$(MIGRATIONS_PATH) -database="$(DB_ADDR)" up

.PHONY: migrate-down
migrate-down:
	@migrate -path=$(MIGRATIONS_PATH) -database="$(DB_ADDR)" down $(filter-out $@,$(MAKECMDGOALS))

.PHONY: migrate-force
migrate-force:
	@migrate -path=$(MIGRATIONS_PATH) -database="$(DB_ADDR)" force $(filter-out $@,$(MAKECMDGOALS))

# -------------------------
# Protobuf
# -------------------------
PROTO_DIR := proto
PROTO_SRC := $(wildcard $(PROTO_DIR)/*.proto)
GO_OUT := .

.PHONY: generate-proto
generate-proto:
	protoc \
		--proto_path=$(PROTO_DIR) \
		--go_out=$(GO_OUT) \
		--go-grpc_out=$(GO_OUT) \
		$(PROTO_SRC)

