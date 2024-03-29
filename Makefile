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

cons-refresh-start: cons-clean-init cons-build-start

cons-build-start: cons-build cons-start

cons-build:
	@docker.local/bin/consumer_build.sh

cons-start:
	@cd docker.local/consumer$(CONSUMER_ID) && ../bin/consumer_start_bls.sh

cons-clean-init:
	@docker.local/bin/consumer_clean.sh
	@docker.local/bin/consumer_init.sh

prov-refresh-start: prov-clean-init prov-build-start

prov-build-start: prov-build prov-start

prov-build:
	@docker.local/bin/provider_build.sh

prov-start:
	@cd docker.local/provider$(PROVIDER_ID) && ../bin/provider_start_bls.sh

prov-clean-init:
	@docker.local/bin/provider_clean.sh
	@docker.local/bin/provider_init.sh

magma-refresh-start: magma-clean-init magma-build-start

magma-build-start: magma-build magma-start

magma-build:
	@docker.local/bin/magma_build.sh

magma-start:
	@cd docker.local/magma && ../bin/magma_start.sh

magma-clean-init:
	@docker.local/bin/magma_clean.sh
	@docker.local/bin/magma_init.sh
