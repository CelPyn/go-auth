# syntax=docker/dockerfile:1

FROM golang:1.17.5-alpine

# Set workdir
WORKDIR /app

# Setup dependencies
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copy source code
COPY main.go ./
COPY pkg ./pkg
COPY resources ./resources

# Build app
RUN go build -o /go-auth

EXPOSE 8080

CMD [ "/go-auth" ]