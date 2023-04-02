.PHONY: migrate
migrate:
	go run github.com/rubenv/sql-migrate/...@latest up

.PHONY: protobuf
protobuf:
	go run github.com/bufbuild/buf/cmd/buf@v1.0.0 generate

.PHONY: protobuf/lint
protobuf/lint:
	buf lint

.PHONY: test
test:
	go test -v -coverpkg=./... ./...