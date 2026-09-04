FROM golang:1.27 AS builder

WORKDIR /app

COPY ./backend/go.mod ./backend/go.sum ./

RUN go mod download 

COPY backend/ ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /docker-gs-ping

RUN CGO_ENABLED=0 GOOS=linux go build -o /migrate ./cmd/migrate

# second stage

FROM scratch

COPY --from=builder /docker-gs-ping /docker-gs-ping

COPY --from=builder /migrate /migrate

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

EXPOSE 8080

ENTRYPOINT [ "/docker-gs-ping" ]
