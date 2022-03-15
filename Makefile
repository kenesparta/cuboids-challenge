
# generate a new migration
# example:
# 	make generate-migration NAME=<your-migration-name>
generate-migration:
	go run cli.go migrate generate $(NAME)

# run all pending migrations
db-migrate:
	go run cli.go migrate up

# rollback last migration
db-migrate-down:
	go run cli.go migrate down

# run all pending migrations in test env
db-migrate-test:
	GO_ENVIRONMENT=test go run cli.go migrate up

# rollback last migration in test env
db-migrate-down-test:
	GO_ENVIRONMENT=test go run cli.go migrate down

# create and migrate db in test and dev env
db-prepare: db-migrate db-migrate-test

# install requirements
install:
	go get ./...
	go get -u -v github.com/spf13/cobra

# prepare the app
prepare: install db-prepare

# run app tests
test:
	go test ./app/tests/...

# run linter
lint:
	golangci-lint run ./...

# start the app
start:
	go run app/main.go
