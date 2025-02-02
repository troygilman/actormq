dev: build
	go run ./cmd/local

ui: build
	go run ./cmd/ui

build:
	protoc --go_out=cluster --proto_path=cluster --go_opt=paths=source_relative -I. cluster/cluster.proto
