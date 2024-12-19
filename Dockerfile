FROM golang:1.23-bullseye as builder

WORKDIR /app

COPY go.* $D/

CMD go mod download

COPY .  /app/

RUN cd /app/ && CGO_ENABLED=0 GOOS=linux go build -v -o server github.com/robinmuhia/helm-charts

# Use the official Alpine image for a lean production container.
# https://hub.docker.com/_/alpine
# https://docs.docker.com/develop/develop-images/multistage-build/#use-multi-stage-builds
FROM alpine:3

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/server /server

CMD ["/server"]
