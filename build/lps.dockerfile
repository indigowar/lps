FROM golang:1.21-alpine AS builder
WORKDIR /lps
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o lps ./cmd/lps

FROM scratch
WORKDIR /lps
COPY --from=builder /lps/lps .
EXPOSE 3000
CMD ["./lps"]