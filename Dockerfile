FROM golang:latest

COPY ./ ./

RUN make install
RUN apt update
RUN apt install -y protoc-gen-go
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
RUN export PATH=$PATH:/go/bin
RUN make gen
RUN make build
RUN make test

CMD ["make", "run"]
