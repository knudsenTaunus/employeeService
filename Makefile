.PHONY: migrate
migrate:
	go run github.com/rubenv/sql-migrate/...@latest up
