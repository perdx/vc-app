
# load env vars
-include .env
export TENANT := $(value TENANT)
export AUTHORITY := $(value AUTHORITY)
export CLIENT_ID := $(value CLIENT_ID)
export CLIENT_SECRET := $(value CLIENT_SECRET)

.PHONY: run
run:
	go run .

.PHONY: test
test:
	go test -v ./...