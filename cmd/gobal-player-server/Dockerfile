FROM golang:1.22-bookworm as builder

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download -x

COPY . .

RUN CGO_ENABLED=0 go build ./cmd/gobal-player-server

FROM gcr.io/distroless/static

COPY --from=builder /build/gobal-player-server .

ENTRYPOINT [ "./gobal-player-server" ]
