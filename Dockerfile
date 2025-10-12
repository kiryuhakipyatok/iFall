FROM golang:1.25-alpine AS builder

WORKDIR /usr/local/src

COPY go.mod go.sum ./
RUN go mod download 

RUN go install -tags='no_postgres no_mysql no_clickhouse no_mssql no_ydb no_vertica no_libsql' github.com/pressly/goose/v3/cmd/goose@latest

COPY . .

# RUN go build -o ./bin/app cmd/app/main.go
RUN --mount=type=cache,target=/root/.cache/go-build \
    go build \
    -trimpath \
    -ldflags="-s -w" \
    -p 4 \
    -o ./bin/app \
    cmd/app/main.go

FROM alpine AS runner

RUN apk add --no-cache tzdata

ENV TZ=Europe/Minsk

RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone
RUN apk add --no-cache make

WORKDIR /app

COPY --from=builder /usr/local/src/bin/app .

COPY --from=builder /go/bin/goose /usr/local/bin/goose

COPY configs /configs

COPY internal/migrations /migrations

EXPOSE 1717

CMD ["./app"]