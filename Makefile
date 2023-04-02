.PHONY: migrate
migrate:
	go run github.com/rubenv/sql-migrate/...@latest up

.PHONY: protobuf
protobuf:
	buf generate

.PHONY: protobuf/lint
protobuf/lint:
	buf lint

