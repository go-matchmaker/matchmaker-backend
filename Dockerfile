# Step 1: Modules caching
FROM golang:1.22.0-alpine3.18 as modules
COPY go.mod go.sum /modules/
WORKDIR /modules
RUN go mod download

# Step 2: Builder
FROM golang:1.22.0-alpine3.18 as builder
LABEL maintainer="bcgocer"
COPY --from=modules /go/pkg /go/pkg
COPY . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/app ./cmd/http

# Step 3: Final
FROM scratch
EXPOSE 8081

# GOPATH for scratch images is /
COPY --from=builder /app/cmd/http/config.yml /
COPY --from=builder /bin/app /app
CMD ["/app"]
