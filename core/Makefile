.PHONY: run
run:
	@go run ./cmd/core --config ./app_config-localhost.yml

.PHONY: test
test:
	@go run gotest.tools/gotestsum

.PHONY: proto
proto:
	@if [ -d "proto" ]; then rm -rf "proto"; fi
	@mkdir "proto"
	@protoc -I ../api/ --go_out=proto --go_opt=paths=source_relative --go-grpc_out=proto --go-grpc_opt=paths=source_relative ../api/core/*.proto
#	@protoc -I ../api/ --go_out=proto --go_opt=paths=source_relative --go-grpc_out=proto --go-grpc_opt=paths=source_relative ../api/core/instrument/v1/*.proto
#	@protoc -I ../api/ --go_out=proto --go_opt=paths=source_relative --go-grpc_out=proto --go-grpc_opt=paths=source_relative ../api/core/account/v1/*.proto
#	@protoc -I ../api/ --go_out=proto --go_opt=paths=source_relative --go-grpc_out=proto --go-grpc_opt=paths=source_relative ../api/core/email/v1/*.proto
#	@protoc -I ../api/ --go_out=proto --go_opt=paths=source_relative --go-grpc_out=proto --go-grpc_opt=paths=source_relative ../api/core/regional/v1/*.proto
#	@protoc -I ../api/ --go_out=proto --go_opt=paths=source_relative --go-grpc_out=proto --go-grpc_opt=paths=source_relative ../api/core/notification/v1/*.proto
#	@protoc -I ../api/ --go_out=proto --go_opt=paths=source_relative --go-grpc_out=proto --go-grpc_opt=paths=source_relative ../api/core/product/v1/*.proto
#	@protoc -I ../api/ --go_out=proto --go_opt=paths=source_relative --go-grpc_out=proto --go-grpc_opt=paths=source_relative ../api/core/registration/v1/*.proto
	@protoc -I ../api/ --go_out=proto --go_opt=paths=source_relative --go-grpc_out=proto --go-grpc_opt=paths=source_relative ../api/core/events/*.proto
	
# REST API
#	@protoc -I ../api/ --grpc-gateway_out ./proto --grpc-gateway_opt logtostderr=true --openapiv2_out ./swagger --openapiv2_opt logtostderr=true --openapiv2_opt use_go_templates=true --grpc-gateway_opt paths=source_relative ../api/core/position/v1/*.proto
#	@protoc -I ../api/ --grpc-gateway_out ./proto --grpc-gateway_opt logtostderr=true --grpc-gateway_opt paths=source_relative ../api/core/account/v1/*.proto
#	@protoc -I ../api/ --grpc-gateway_out ./proto --grpc-gateway_opt logtostderr=true --grpc-gateway_opt paths=source_relative ../api/core/product/v1/*.proto
	@protoc -I ../api/ --grpc-gateway_out ./proto --grpc-gateway_opt logtostderr=true --grpc-gateway_opt paths=source_relative ../api/core/events/*.proto


.PHONY: build
build:
	@del /Q core.exe
	@go clean
	@go build -race -v ./cmd/core
