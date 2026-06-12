FROM golang:1.26-alpine AS builder

WORKDIR /app
COPY go.mod go.sum* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /app/server ./cmd/server

FROM alpine:3.21
WORKDIR /app
COPY --from=builder /app/server .
COPY data/ ./data/
COPY index/ ./index/
EXPOSE 8080
CMD ["./server"]
