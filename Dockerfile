FROM golang:alpine AS build-env
LABEL version="1.0"
LABEL maintainer="abohmeed@gmail.com"
RUN mkdir /go/src/app && apk update && apk add git && go get -u github.com/golang/dep/cmd/dep
ADD src/github.com/abohmeed/* /go/src/app/
COPY ./Gopkg.toml /go/src/app
WORKDIR /go/src/app 
RUN dep ensure && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o main .

FROM scratch
WORKDIR /app
COPY --from=build-env /go/src/app/main /app/
ENTRYPOINT [ "./main" ]