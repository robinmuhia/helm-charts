FROM golang:1.23-bullseye as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY .  ./

RUN CGO_ENABLED=0 GOOS=linux go build -v -o server github.com/robinmuhia/helm-charts

# Use the official Alpine image for a lean production container.
# https://hub.docker.com/_/alpine
# https://docs.docker.com/develop/develop-images/multistage-build/#use-multi-stage-builds
FROM alpine:3

RUN apk add --no-cache ca-certificates

# install helm
RUN apk add --no-cache curl bash && \
    curl -fsSL https://get.helm.sh/helm-v3.11.2-linux-amd64.tar.gz -o helm.tar.gz && \
    tar -zxvf helm.tar.gz && \
    mv linux-amd64/helm /usr/local/bin/helm && \
    rm -rf helm.tar.gz linux-amd64

COPY --from=builder /app/server /server

CMD ["/server"]
