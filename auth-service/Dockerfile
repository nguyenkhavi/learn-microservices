FROM golang:1.19-alpine as builder
RUN apk --no-cache add ca-certificates git
WORKDIR /build

# Fetch dependencies
COPY go.mod go.sum ./
RUN go mod download

# Build
COPY . ./
RUN CGO_ENABLED=0 go build -o /build/myapp .

# Create final image
FROM scratch as production
WORKDIR /
COPY --from=builder /build/myapp .
COPY --from=builder /build/.env .env

EXPOSE 8080
CMD ["./myapp"]