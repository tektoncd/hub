FROM golang:1.15-alpine3.12 AS builder

WORKDIR /go/src/github.com/tektoncd/hub/api
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o api-server ./cmd/api/...

FROM alpine:3.12

RUN apk --no-cache add git ca-certificates && addgroup -S hub && adduser -S hub -G hub
USER hub

WORKDIR /app
COPY --from=builder /go/src/github.com/tektoncd/hub/api/api-server /app/api-server
COPY --from=builder /go/src/github.com/tektoncd/hub/api/gen/http/openapi3.yaml /app/gen/http/openapi3.yaml
EXPOSE 8000

CMD [ "/app/api-server" ]