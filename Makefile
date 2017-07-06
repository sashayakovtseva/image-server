bin:
	go build  --ldflags '-linkmode "external" -extldflags "-static"' -x -o image-server .
build:
	sudo docker build -t sashayakovtseva/iotracks:image-server$(TAG) .
push:build
	sudo docker push sashayakovtseva/iotracks:image-server$(TAG)
