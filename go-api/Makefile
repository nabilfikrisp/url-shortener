.PHONY: test test-unit test-integration coverage db-up db-down wait-for-db

# Run all tests
test:
	go test ./... -cover

# Run only unit tests
test-unit:
	go test ./tests/unit/... -v

# Run integration tests
test-integration:
	go test ./tests/integration/... -v

# Bring up docker compose test db
test-db-up:
	docker compose -f docker-compose.test.yml up -d test-db

# Tear down docker compose test db
test-db-down:
	docker compose -f docker-compose.test.yml down -v

coverage:
	go test ./... -coverpkg=./... -coverprofile=coverage.out
	go tool cover -html=coverage.out