FROM golang:1.18.2-alpine3.14 AS builder

WORKDIR /go/src/github.com/tektoncd/hub
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o api-server ./api/cmd/api/...

FROM alpine:3.16.0

RUN apk --no-cache add git ca-certificates openssh-client && addgroup -S hub && adduser -S hub -G hub
USER hub

WORKDIR /app

COPY --from=builder /go/src/github.com/tektoncd/hub/api-server /app/api-server

# For each new version, doc has to be copied
COPY api/gen/http/openapi3.json /app/docs/openapi3.json
COPY api/v1/gen/http/openapi3.json /app/docs/v1/openapi3.json

EXPOSE 8000

CMD [ "/app/api-server" ]
