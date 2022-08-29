# Stage 1: build the target binaries
FROM golang:1.18 AS builder

WORKDIR /workspace

COPY go.mod go.sum ./
COPY ./lifecycle/ ./lifecycle/

RUN go build -o bin/lifecycle ./lifecycle/main.go

# Stage 2: Copy the binaries from the image builder to the base image
FROM alpine:latest  

COPY --from=builder /workspace/bin/lifecycle /bin/lifecycle
COPY ./lifecycle/hooks/poststart.sh /bin/poststart.sh
COPY ./lifecycle/hooks/prestop.sh /bin/prestop.sh

ENTRYPOINT ["/bin/lifecycle"]