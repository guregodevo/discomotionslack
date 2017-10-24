FROM golang:alpine

EXPOSE 8080

RUN mkdir -p /var/log/discomotion

WORKDIR /go/src/github.com/dailymotion-leo/discomotion
COPY . . 

RUN go install 

ENTRYPOINT /go/bin/discomotion --config dev.config.yaml 


