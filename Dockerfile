FROM golang:1.16-alpine as builder

RUN apk update && \
    apk add --no-cache git

WORKDIR /gosuslugi

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

ENV CGO_ENABLED 0
RUN go build -o bin/server cmd/server/*.go
RUN go build -o bin/seed cmd/seed/*.go
EXPOSE 3000

FROM scratch
COPY --from=builder /gosuslugi/bin .
CMD ["./server"]
