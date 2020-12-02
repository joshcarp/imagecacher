FROM golang:1.14-alpine as installer
WORKDIR /usr/app
ADD . .
RUN go build -o imagecacher .
FROM alpine:latest
WORKDIR /usr/app
COPY --from=installer /usr/app/imagecacher /bin/imagecacher
ENTRYPOINT ["./imagecacher"]