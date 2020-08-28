FROM golang:latest as builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPROXY=https://goproxy.cn,direct

WORKDIR /src

COPY ./ /src

RUN go mod download

RUN go build -o dist/main


FROM scratch

COPY --from=builder /src/dist/main /

# Command to run
ENTRYPOINT ["/main", "serve"]