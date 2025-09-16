# Use latest Go - adjust as needed
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /go-books

FROM scratch
COPY --from=builder /go-books /go-books
EXPOSE 8080
ENTRYPOINT ["/go-books"]
