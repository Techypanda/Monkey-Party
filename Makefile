# https://gist.github.com/turtlemonvh/38bd3d73e61769767c35931d8c70ccb4

compile-grpc: # Compiles GRPC
	protoc --go_out=pkg/rooms_grpc/ --go_opt=paths=source_relative --go-grpc_out=pkg/rooms_grpc/ --go-grpc_opt=paths=source_relative rooms.proto

test: # TODO
	go test -v ./... -coverprofile coverprofile

roomagent:
	go run cmd/roomagent/main.go

monkeypartyapi:
	go run cmd/monkeypartyapi/main.go

grpc-gui: # Uses https://github.com/fullstorydev/grpcui to give you a gui for roomagent
	grpcui -plaintext localhost:8054

test-roomagent:
	go test -v --cover ./cmd/roomagent

clean: # TODO

.PHONY: test clean roomagent test-roomagent compile-grpc