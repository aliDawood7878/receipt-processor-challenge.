FROM golang:1.20 AS builder

WORKDIR /app

# Copy go.mod first
COPY go.mod ./

COPY main.go ./

RUN go mod tidy

# Build a static binary
ENV CGO_ENABLED=0
ENV GOOS=linux
RUN go build -o receipt-processor main.go

FROM gcr.io/distroless/static:latest
COPY --from=builder /app/receipt-processor /receipt-processor

EXPOSE 8080
USER nonroot:nonroot
CMD ["/receipt-processor"]