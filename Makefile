dev: build
	go run ./cmd/local

build:
	protoc --go_out=raft --proto_path=raft --go_opt=paths=source_relative -I. raft/raft.proto
