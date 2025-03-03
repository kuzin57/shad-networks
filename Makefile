gen:
	protoc -I proto proto/graph.proto --go_out=internal/generated --go_opt=paths=source_relative --go-grpc_out=internal/generated --go-grpc_opt=paths=source_relative

clean:
	rm ./internal/generated/*.pb.go && rm ./main

install:
	go get -d ./...

test:
	go test -v ./...

build:
	go build ./cmd/main

run:
	./main --config ./config/config.yaml
