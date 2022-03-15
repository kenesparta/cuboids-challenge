
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

# Generate the HTML report, needs: Needs pip install junit2html
test-html:
	go test ./app/tests/... -v ;\
	rm -rf reports/controllers.html ;\
	rm -rf reports/model.html ;\
	junit2html reports/controllers.xml reports/controllers.html ;\
	junit2html reports/model.xml reports/model.html

# Generate the HTML report, needs: Needs pip install junit2html
lint-html:
	rm -rf reports/lint.html
	golangci-lint run --out-format junit-xml ./... > reports/lint.xml ;\
	junit2html reports/lint.xml reports/lint.html
