FROM golang:1.24-alpine3.21 as base

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -v -o ./app ./cmd/main.go

FROM alpine:3.21 as final

WORKDIR /app

COPY --from=base /build/app /app/app

CMD ["/app/app"]

