FROM golang:latest

COPY ./ ./

RUN make install
RUN apt update
RUN apt install -y protoc-gen-go
RUN apt install -y graphviz
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
RUN export PATH=$PATH:/go/bin
RUN make gen
RUN make build-graph

CMD ["make", "run-graph"]
