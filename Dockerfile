FROM golang:1.11

WORKDIR /image-server
COPY main.go .
RUN go build --ldflags '-linkmode "external" -extldflags "-static"' -o image-server main.go

FROM scratch
COPY images/ images
COPY --from=0 /image-server/image-server .
EXPOSE 8080/tcp
CMD ["./image-server"]
