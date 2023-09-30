# https://gist.github.com/turtlemonvh/38bd3d73e61769767c35931d8c70ccb4

test: # TODO
	go test -v --cover ./...

roomagent:
	go run cmd/roomagent/main.go

test-roomagent:
	go test -v --cover ./cmd/roomagent

clean: # TODO

.PHONY: test clean roomagent test-roomagent