FROM golang:1.18.2 as builder

WORKDIR /src
COPY . /src

RUN go build -buildmode=c-shared -o out_rabbitmq.so .

FROM fluent/fluent-bit:1.9.3 as fluent-bit
USER root

COPY --from=builder /src/out_rabbitmq.so /fluent-bit/bin/
COPY --from=builder /src/out_rabbitmq.h /fluent-bit/bin/
COPY --from=builder /src/*.conf /fluent-bit/etc/
