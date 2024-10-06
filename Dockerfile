FROM golang:1.23-alpine AS builder
LABEL authors="D0niL"

WORKDIR /usr/local/src

RUN apk --no-cache add bash git make gcc gettext musl-dev

COPY ["go.mod", "go.sum", "./"]
RUN go mod download

COPY . .
RUN go build -o ./bin/app ./cmd/clinic/main.go

FROM alpine AS runner

COPY --from=builder /usr/local/src/bin/app /
COPY --from=builder /usr/local/src/internal/storage/postgres/migrations /migrations
COPY config/prod.yaml /prod.yaml
ENV CONFIG_PATH=/prod.yaml
EXPOSE 8080

CMD ["/app"]
