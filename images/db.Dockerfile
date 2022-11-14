FROM --platform=$BUILDPLATFORM golang:1.18.2-alpine3.14 AS builder

WORKDIR /go/src/github.com/tektoncd/hub
COPY . .
ARG TARGETOS
ARG TARGETARCH

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o db-migration ./api/cmd/db/...

FROM alpine:3.16.3

RUN apk --no-cache add ca-certificates && addgroup -S hub && adduser -S hub -G hub
USER hub

WORKDIR /app
COPY --from=builder /go/src/github.com/tektoncd/hub/db-migration /app/db-migration
EXPOSE 8000

CMD [ "/app/db-migration" ]
