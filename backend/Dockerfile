# Homebox builder
FROM golang:latest AS builder
LABEL Maintainer="lans.rf@gmail.com" 


RUN apt update -y && \
    apt upgrade -y && \
    mkdir -p /homebox/

COPY endpoint /homebox/endpoint
COPY main.go /homebox/main.go
COPY go.mod /homebox/go.mod

RUN cd /homebox/ && go build .

# Homebox image
FROM alpine:latest
LABEL Maintainer="lans.rf@gmail.com"

RUN apk add gcompat && \
    mkdir /homebox/

COPY --from=builder /homebox/homebox /homebox/homebox
ENTRYPOINT ["/homebox/homebox"]
