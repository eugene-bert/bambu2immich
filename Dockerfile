FROM golang:1.27-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /bambu2immich .

FROM alpine:3.21
RUN apk add --no-cache ca-certificates \
    && adduser -D -u 1000 -H appuser \
    && mkdir -p /data/timelapses \
    && chown -R appuser:appuser /data
COPY --from=build /bambu2immich /bambu2immich
USER appuser
ENTRYPOINT ["/bambu2immich"]
