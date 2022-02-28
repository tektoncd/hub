FROM golang:1.15-alpine3.14 AS builder

WORKDIR /go/src/github.com/tektoncd/hub
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o db-migration ./api/cmd/db/...

FROM alpine:3.14

RUN apk --no-cache add ca-certificates && addgroup -S hub && adduser -S hub -G hub
USER hub

WORKDIR /app
COPY --from=builder /go/src/github.com/tektoncd/hub/db-migration /app/db-migration
EXPOSE 8000

CMD [ "/app/db-migration" ]
