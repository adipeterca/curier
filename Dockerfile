FROM golang:1.26 AS builder

WORKDIR /app

# No additional dependencies, no go.sum
# COPY go.mod go.sum ./
COPY go.mod ./
RUN go mod download

COPY static/ ./static/
COPY templates/ ./templates/
COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o curier .

FROM scratch

WORKDIR /

ENV CURIER_HOST="0.0.0.0"
ENV CURIER_PORT="39800"
ENV CURIER_STORAGE_PATH="/uploads/" 

COPY --from=builder /app/curier /curier

EXPOSE 39800

ENTRYPOINT ["/curier"]