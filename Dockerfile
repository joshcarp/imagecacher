FROM golang:1.14-alpine as installer
ADD . .
RUN go install github.com/joshcarp/imagecacher
WORKDIR /usr/app
ENTRYPOINT ["imagecacher"]