FROM golang:1.14-buster
WORKDIR /usr/app
ADD . .
RUN go install github.com/joshcarp/imagecacher
ENTRYPOINT ["imagecacher"]