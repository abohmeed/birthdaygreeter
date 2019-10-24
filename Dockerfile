FROM golang:alpine AS build-env
ADD . /
RUN export GOPATH=/ && cd /src/github.com/abohmeed/birthdaygreeter && go install

FROM alpine
WORKDIR /app
COPY --from=build-env /bin/birthdaygreeter /app/
ENTRYPOINT [ "./birthdaygreeter" ]