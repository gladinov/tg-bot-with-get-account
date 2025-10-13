FROM golang:1.25.1-alpine AS builder


WORKDIR /usr/local/src

RUN apk add --no-cache bash git make gettext gcc musl-dev

# dependencies
COPY ["go.mod","go.sum","./"]
RUN go mod download

# # build 
COPY . ./
RUN CGO_ENABLED=1 GOOS=linux go build -o ./bin/app main.go

FROM alpine AS runner

COPY --from=builder /usr/local/src/bin/app /

EXPOSE 8080


CMD ["/app"]