folder: ## create folder structure
	@mkdir -p configs
	@mkdir -p registry
	@mkdir -p db
	@mkdir -p log
	@mkdir -p util
	@mkdir -p cache
	@mkdir -p opentracing
	@mkdir -p opentracing/jaeger
	@mkdir -p interceptor

dep: ## Install the dependencies
	@go get \
		github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
		github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
		github.com/gogo/protobuf/protoc-gen-gogo \
		github.com/mwitkow/go-proto-validators/protoc-gen-govalidators \
		google.golang.org/protobuf/cmd/protoc-gen-go \
		google.golang.org/grpc/cmd/protoc-gen-go-grpc \

mod: ## Install go library
	@go mod tidy
	@go mod vendor

genpb:
	@protoc \
	-I api/proto \
	-I third_party \
	-I vendor/github.com/grpc-ecosystem/grpc-gateway/v2 \
	-I vendor/github.com/mwitkow/go-proto-validators \
	-I vendor/github.com/gogo/protobuf \
	--go_out=api/pb \
	--go-grpc_out=api/pb \
	--govalidators_out=gogoimport=true:api/pb \
	--grpc-gateway_out=api/pb \
    --openapiv2_out=docs \
	api/proto/*.proto

run:
	@statik -f -src=docs -dest=cmd/server && cd cmd/server && go run main.go