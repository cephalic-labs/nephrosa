FROM golang:1.26.10-alpine AS build

RUN apk add --no-cache ca-certificates
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY backend/ ./backend/
COPY db/ ./db/
RUN CGO_ENABLED=0 go build \
    -trimpath \
    -ldflags="-s -w -buildid=" \
    -o /out/renal-relay \
    ./backend/cmd/api

FROM scratch
COPY --from=build /out/renal-relay /renal-relay
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
USER 65532:65532
EXPOSE 8080
ENTRYPOINT ["/renal-relay"]
