gen:
	protoc -I proto proto/*.proto --go_out=./generated --go_opt=paths=source_relative --go-grpc_out=./generated --go-grpc_opt=paths=source_relative

clean:
	rm -f ./generated/*.pb.go
	rm -f ./main

install:
	go get ./...

test-%:
	go test -v ./services/$*/...

build-%:
	go build ./services/$*/cmd/main

run-%:
	./main --config ./services/$*/config/config.yaml --secrets ./services/$*/config/secrets.yaml

.PHONY: graph
