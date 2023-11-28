ARG GO_IMAGE
ARG GO_IMAGE_VERSION
ARG GO_VER
ARG ALPINE_IMAGE
ARG ALPINE_VER

FROM ${GO_IMAGE}:${GO_VER}-alpine as builder

RUN apk --no-cache update && \
    apk --no-cache add git gcc libc-dev

WORKDIR /opt/app

COPY . .

ENV CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64 \
    GO111MODULE=on

RUN go mod download

RUN go build -tags musl -o service ./main.go

FROM ${ALPINE_IMAGE}:${ALPINE_VER}
WORKDIR /opt/app
COPY --from=builder /opt/app/service /opt/app/service 

USER daemon

EXPOSE 8000

CMD ["./service"]
