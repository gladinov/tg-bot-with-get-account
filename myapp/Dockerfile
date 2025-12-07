FROM golang:1.25.1-alpine AS builder

WORKDIR /usr/local/src

RUN apk add --no-cache bash git make gettext gcc musl-dev

# dependencies
COPY ["go.mod","go.sum","./"]
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -o ./bin/myapp ./cmd/main.go

FROM alpine AS runner

#RUN mkdir -p /configs

COPY --from=builder /usr/local/src/bin/myapp /

#COPY --from=builder /usr/local/src/configs/config.yaml /configs/
#COPY --from=builder /usr/local/src/configs/sber.yaml /configs/

EXPOSE 8080


CMD ["/myapp"]