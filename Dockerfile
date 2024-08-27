# Start from golang base image
FROM golang:1.21.6 AS builder

# Set the current working directory inside the container 
WORKDIR /app

ENV GOBIN=/app/bin

# Copy go mod and sum files 
COPY go.mod go.sum ./


# Copy the source from the current directory to the working Directory inside the container 
COPY ./cmd cmd
COPY ./src src

RUN --mount=type=cache,target=/go --mount=type=cache,target=/root/.cache go install -race ./cmd/...

FROM debian

COPY --from=builder /app/bin/ /app/

RUN apt-get update && apt-get install -y ca-certificates

ENV PATH=/app:$PATH

WORKDIR /app

CMD ["service"]
