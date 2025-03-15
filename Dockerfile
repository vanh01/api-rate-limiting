FROM golang:1.23-alpine as go-builder

# Add git package
RUN apk add --no-cache git
WORKDIR /app_service

COPY . .

RUN go mod tidy && \
    go build -o ./server ./cmd

FROM alpine:latest as arl

WORKDIR /app_service

RUN apk --no-cache add ca-certificates

COPY --from=go-builder /app_service/server ./server
COPY --from=go-builder /app_service/config.yaml ./config.yaml

EXPOSE 8080

ENTRYPOINT [ "/app_service/server" ]