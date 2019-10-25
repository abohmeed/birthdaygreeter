FROM golang:alpine AS build-env
LABEL version="1.0"
LABEL maintainer="abohmeed@gmail.com"
RUN mkdir /go/src/app /go/src/redis-check && apk update && apk add git && go get -u github.com/golang/dep/cmd/dep
ADD src/github.com/abohmeed/birthdaygreeter/* /go/src/app/
ADD src/github.com/abohmeed/redis-check/* /go/src/redis-check/
WORKDIR /go/src/birthdaygreeter
RUN dep ensure && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o app .
WORKDIR /go/src/redis-check
RUN dep ensure && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o redis-check .

FROM scratch
WORKDIR /app
COPY --from=build-env /go/src/birthdaygreeter/app /app/
COPY --from=build-env /go/src/redis-check/redis-check /app/
ENTRYPOINT [ "./app" ]