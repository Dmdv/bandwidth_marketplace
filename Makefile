# Makefile
# Specific features to manage automations in development & build processes.
#
# For usage on Windows see [Chocolatey CLI Documentation](https://docs.chocolatey.org/en-us/choco/setup)
# Then execute `choco install make` command in shell, now you will be able to use `make` on Windows.

make_path := $(abspath $(lastword $(MAKEFILE_LIST)))
root_path := $(patsubst %/, %, $(dir $(make_path)))

.PHONY: pre-push go-mod check-commit
pre-push: go-mod check-commit

.PHONY: check-commit go-get run-test run-lint
check-commit: go-get generate run-test run-lint

go-mod:
	@echo "Prepare Go mod files..."

	@cd $(root_path)/code/consumer && go mod tidy -v
	@cd $(root_path)/code/consumer && go mod download -x
	@cd $(root_path)/code/pb && go mod tidy -v
	@cd $(root_path)/code/pb && go mod download -x
	@cd $(root_path)/code/magma && go mod tidy -v
	@cd $(root_path)/code/magma && go mod download -x
	@cd $(root_path)/code/core && go mod tidy -v
	@cd $(root_path)/code/core && go mod download -x
	@cd $(root_path)/code/provider && go mod tidy -v
	@cd $(root_path)/code/provider && go mod download -x

	@echo "Go mod files completed."

go-get:
	@echo "Load dependencies..."

	@cd $(root_path)/code/consumer && go get -d ./...
	@cd $(root_path)/code/pb && go get -d ./...
	@cd $(root_path)/code/magma && go get -d ./...
	@cd $(root_path)/code/core && go get -d ./...
	@cd $(root_path)/code/provider && go get -d ./...

	@echo "Dependencies loaded."

run-test:
	@echo "Start testing..."

	@cd $(root_path)/code/consumer && go test -tags bn256 -cover ./...
	@cd $(root_path)/code/provider && go test -tags bn256 -cover ./...

	@echo "Tests completed."

run-lint:
	@echo "Start linters..."

	@echo "Checking consumer module..."
	@cd $(root_path)/code/consumer && golangci-lint run

	@echo "Checking magma module..."
	@cd $(root_path)/code/magma && golangci-lint run

	@echo "Checking provider module..."
	@cd $(root_path)/code/provider && golangci-lint run

	@echo "Linters completed."

proto_path=./code

generate:
	@echo "Compiling protobuf files..."

	@protoc -I $(proto_path) \
	--go_opt=module="github.com/0chain/bandwidth_marketplace/code" \
	--go-grpc_opt=module="github.com/0chain/bandwidth_marketplace/code" \
	--go-grpc_out=$(proto_path) \
	--go_out=$(proto_path) \
	--grpc-gateway_out=$(proto_path) \
	$(proto_path)/pb/consumer/proto/*.proto

	@protoc -I $(proto_path) \
	--go_opt=module="github.com/0chain/bandwidth_marketplace/code" \
    --go-grpc_opt=module="github.com/0chain/bandwidth_marketplace/code" \
	--go-grpc_out=$(proto_path) \
	--go_out=$(proto_path) \
	--grpc-gateway_out=$(proto_path) \
	$(proto_path)/pb/provider/proto/*.proto

	@echo "Compiling completed."
