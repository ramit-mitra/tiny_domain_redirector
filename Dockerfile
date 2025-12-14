# build stage
FROM golang:latest AS builder

WORKDIR /app

# copy go.mod and go.sum files to download dependencies
COPY go.mod go.sum ./
RUN go mod download

# copy the rest of the source code
COPY . .

# build the app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/tiny_domain_redirector .

# final stage
FROM alpine:latest

WORKDIR /app

# copy the built binary from the builder stage
COPY --from=builder /app/tiny_domain_redirector .

# copy the redirects file
COPY redirects.yaml .

# expose the application
EXPOSE 9990

# run the app
CMD ["/app/tiny_domain_redirector"]
