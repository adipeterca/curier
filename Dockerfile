FROM golang:1.26 AS builder

WORKDIR /app

# No additional dependencies, no go.sum
# COPY go.mod go.sum ./
COPY go.mod ./
RUN go mod download

COPY static/ ./static/
COPY templates/ ./templates/
COPY *.go ./

RUN mkdir -p /uploads

ARG VERSION=dev
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-X main.version=${VERSION}" -o curier .

FROM scratch

WORKDIR /

ENV CURIER_HOST="0.0.0.0"
ENV CURIER_PORT="39800"
ENV CURIER_STORAGE_PATH="/uploads/" 

COPY --from=builder /app/curier /curier
COPY --from=builder /tmp /tmp
COPY --from=builder /uploads /uploads

EXPOSE 39800

ENTRYPOINT ["/curier"]