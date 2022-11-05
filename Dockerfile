FROM node:10.15.3 as NODE_BUILD
WORKDIR /go/src/github.com/leimeng-go/pipe/
ADD . /go/src/github.com/leimeng-go/pipe/
RUN cd console && npm install && npm run build && cd ../theme && npm install && npm run build && \
    rm -rf node_modules && cd ../console && rm -rf node_modules

FROM golang:alpine as GO_BUILD
WORKDIR /go/src/github.com/leimeng-go/pipe/
COPY --from=NODE_BUILD /go/src/github.com/leimeng-go/pipe/ /go/src/github.com/leimeng-go/pipe/
ENV GO111MODULE=on
RUN apk add --no-cache gcc musl-dev git && go build -i -v

FROM alpine:latest
LABEL maintainer="Liang Ding<845765@qq.com>"
WORKDIR /opt/pipe/
COPY --from=GO_BUILD /go/src/github.com/leimeng-go/pipe/ /opt/pipe/
RUN apk add --no-cache ca-certificates tzdata

ENV TZ=Asia/Shanghai
EXPOSE 5897

ENTRYPOINT [ "/opt/pipe/pipe" ]
