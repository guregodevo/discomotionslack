FROM golang:alpine

EXPOSE 8080

RUN mkdir -p /var/log/discomotionslack

WORKDIR /go/src/github.com/dailymotion-leo/discomotionslack
COPY . . 

RUN go install 

ENTRYPOINT /go/bin/discomotionslack --config dev.config.yaml 


