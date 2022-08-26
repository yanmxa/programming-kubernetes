# Stage 1: build the target binaries
FROM golang:1.18 AS builder

WORKDIR /workspace

COPY go.mod ./
COPY ./lifecycle/ ./lifecycle/

RUN go build -o bin/lifecycle ./lifecycle/main.go

# Stage 2: Copy the binaries from the image builder to the base image
FROM alpine:latest  

COPY --from=builder /workspace/bin/lifecycle /usr/local/bin/lifecycle
COPY ./lifecycle/hooks/prestop.sh /usr/local/bin/prestop.sh
COPY ./lifecycle/hooks/poststart.sh /usr/local/bin/poststart.sh

ENTRYPOINT ["/usr/local/bin/lifecycle"]