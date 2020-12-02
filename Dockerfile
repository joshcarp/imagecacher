FROM golang:1.14-alpine as installer
WORKDIR /usr/app
ADD . .
RUN go build -o imagecacher .
FROM alpine:latest
COPY --from=installer /usr/app/imagecacher /bin/imagecacher
ENTRYPOINT ["./imagecacher"]