FROM golang:1.24 AS builder

WORKDIR /app

# No additional dependencies, no go.sum
# COPY go.mod go.sum ./
COPY go.mod ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o curier .

FROM scratch

COPY --from=builder /app/curier /curier

EXPOSE 8080

ENTRYPOINT ["/curier"]