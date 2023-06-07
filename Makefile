.PHONY:proto
proto:
	protoc --go-grpc_out=. --go_out=. proto/api.proto

image:
	KO_DOCKER_REPO=hongliang5316 ko build ./cmd/midjourney-apiserver -B --platform=all
