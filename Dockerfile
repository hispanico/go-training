FROM golang:1.21.3-alpine AS builder

COPY ./src /go/src/app/

WORKDIR /go/src/app/

RUN mkdir bin && go build -o bin/webserver ./cmd/server/main.go

FROM gcr.io/distroless/static-debian11

COPY --from=builder /go/src/app/bin/webserver /app/webserver

CMD ["/app/webserver"]
