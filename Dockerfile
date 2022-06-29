FROM golang:1.17-buster as builder

WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o application ./cmd/main.go

FROM alpine:3.15.4
ENV PG_URL="postgres://postgres:1111@mts_teta_projects-db-1:5432/mtsteta"
ENV GRPC_URL="mts_teta_projects-auth-1:4000"
COPY --from=builder /app/application /app/application
CMD ["/app/application"]
