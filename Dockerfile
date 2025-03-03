FROM golang:latest

COPY ./ ./

RUN make install

RUN make build

RUN make test

CMD ["make", "run"]
