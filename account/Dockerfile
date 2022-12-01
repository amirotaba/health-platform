# syntax=docker/dockerfile:1
FROM golang:1.17.5-alpine3.14 AS development
ENV PROJECT_PATH=/account
ENV PATH=$PATH:$PROJECT_PATH/build
ENV CGO_ENABLED=0
ENV GO_EXTRA_BUILD_ARGS="-a -installsuffix cgo"

# RUN apk add --no-cache ca-certificates tzdata make git bash protobuf
RUN --mount=type=cache,mode=0777,id=apkcache,target=/var/cache/apk \
    --mount=type=cache,mode=0777,id=apklib,target=/var/lib/apk \
    apk add ca-certificates tzdata make git bash protobuf

RUN mkdir -p $PROJECT_PATH
COPY . $PROJECT_PATH
WORKDIR $PROJECT_PATH

ENV GOPROXY=https://goproxy.io,direct

ENV CGO_ENABLED=0
#COPY go.* .
RUN #go mod download

#COPY  ./config.yaml /account/build/config.yaml
COPY  ./setup.sh /account/build/setup.sh
#COPY  internal/app/database/migrate/migrations /account/build/migrations
#ARG TARGETOS
#ARG TARGETARCH
#RUN go mod download
#RUN --mount=type=cache,target=/root/.cache/go-build \
#GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o build/account ./main.go
# RUN go build $(GO_EXTRA_BUILD_ARGS) -ldflags "-s -w -X main.version=$(VERSION)" -o build/account ./main.go
RUN --mount=type=cache,mode=0777,id=gocache,target=/root/.cache/go-build \
    --mount=type=cache,mode=0777,id=gomodcache,target=/go/pkg/mod \
      go build $(GO_EXTRA_BUILD_ARGS) -ldflags "-s -w -X main.version=$(VERSION)" -o build/account ./main.go
FROM alpine:3.15.0 AS production

RUN apk --no-cache add ca-certificates tzdata
#COPY --from=development /account/build/config.yaml /config.yaml
COPY --from=development /account/build/setup.sh /setup.sh
#COPY --from=development /account/build/migrations /migrations
COPY --from=development /account/build/account /usr/bin/account
#USER nobody:nogroup
EXPOSE 17070
RUN #chmod 777 setup.sh
RUN ["chmod", "+x", "/setup.sh"]
ENTRYPOINT ["/setup.sh"]
#FROM golang:1.14 AS builder
#
#ENV APP_HOME /app
#
#WORKDIR $APP_HOME
#
#ENV GOBIN /go/bin
#RUN go env -w GO111MODULE=on
#RUN go env -w CGO_ENABLED=0
#RUN go env -w GOPROXY=https://goproxy.io,direct
#
#COPY go.mod go.sum ./
#RUN go mod tidy
#
#RUN mkdir -p $APP_HOME/static
#
#ADD ./internal $APP_HOME/internal
#ADD ./main.go $APP_HOME
#ADD ./config.yaml $APP_HOME
#RUN go build -v -o /go/bin/{$APP_HOME}
#
#FROM alpine:latest AS app
#
#ARG PROJECT
#
#COPY --from=builder /go/bin/${PROJECT} /
#
#RUN chmod +x /${PROJECT}
#
#ENTRYPOINT /${PROJECT}
