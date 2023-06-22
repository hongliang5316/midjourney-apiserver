version=0.0.2

proto:
	protoc --proto_path=./pkg/api --go-grpc_out=. --go_out=. pkg/api/api.proto
	protoc --proto_path=./pkg/api --go-grpc_out=. --go_out=. pkg/api/common.proto
	protoc --proto_path=./pkg/api --go-grpc_out=. --go_out=. pkg/api/imagine.proto
	protoc --proto_path=./pkg/api --go-grpc_out=. --go_out=. pkg/api/upscale.proto
	protoc --proto_path=./pkg/api --go-grpc_out=. --go_out=. pkg/api/describe.proto

image:
	KO_DOCKER_REPO=hongliang5316 ko build ./cmd/midjourney-apiserver -B --platform=all -t $(version)
