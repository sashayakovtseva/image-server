FROM alpine:latest
#FROM hypriot/rpi-alpine-scratch

COPY image-server /go/bin/
COPY images /go/bin/images/
WORKDIR /go/bin
CMD ["./image-server"]
