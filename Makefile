version=0.0.2

proto:
	protoc --go-grpc_out=. --go_out=. pkg/api/api.proto

image:
	KO_DOCKER_REPO=hongliang5316 ko build ./cmd/midjourney-apiserver -B --platform=all -t $(version)
