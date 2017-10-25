FROM golang:alpine

EXPOSE 8080

RUN mkdir -p /var/log/discomotionslack

WORKDIR /go/src/github.com/guregodevo/discomotionslack
COPY . . 

RUN go install 

ENTRYPOINT /go/bin/discomotionslack --config config/dev.config.yaml 


