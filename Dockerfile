FROM golang:1.27 AS builder

WORKDIR /app

COPY ./backend/go.mod ./backend/go.sum ./

RUN go mod download 

COPY backend/ ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /docker-gs-ping

# second stage

FROM alpine:latest

WORKDIR /

COPY --from=builder /docker-gs-ping /docker-gs-ping

EXPOSE 8080

CMD [ "/docker-gs-ping" ]
