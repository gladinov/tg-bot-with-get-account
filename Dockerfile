FROM golang:1.25.1-alpine AS builder


WORKDIR /usr/local/src

RUN apk add --no-cache bash git make gettext gcc musl-dev

# dependencies
COPY ["myapp/go.mod","myapp/go.sum","./"]
RUN go mod download

# # build 
COPY myapp ./
RUN CGO_ENABLED=1 GOOS=linux go build -o ./bin/myapp main.go

FROM alpine AS runner

RUN mkdir -p /configs

COPY --from=builder /usr/local/src/bin/myapp /

COPY --from=builder /usr/local/src/configs/tinkoffApiConfig.yaml /configs/
COPY --from=builder /usr/local/src/configs/sber.yaml /configs/

EXPOSE 8080


CMD ["/myapp"]