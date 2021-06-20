FROM golang:1.16-alpine as builder

RUN apk update && \
    apk add --no-cache git

WORKDIR /gosuslugi

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o bin/gosuslugi cmd/server/*.go
EXPOSE 3000

FROM scratch
COPY --from=builder /gosuslugi/bin/gosuslugi .
ENTRYPOINT ["./gosuslugi"]
