
TARGET=news_watch_notice
PKG=$(TARGET)
TAG=latest
IMAGE_PREFIX?=hanyajun
IMAGE_PREFIX_PRD=hanyajun
TARGET_IMAGE_DEV=$(IMAGE_PREFIX)/$(TARGET):$(TAG)
TARGET_IMAGE_PRD=$(IMAGE_PREFIX_PRD)/$(TARGET):$(TAG)
all: image

$(TARGET):
	CGO_ENABLED=0 go build -o dist/$(TARGET) $(PKG)/cmd

gitlog:


target:
	mkdir -p dist
	git log | head -n 1 > dist/news_watch_notice.sha
	docker run --rm -i -v `pwd`:/go/src/$(PKG) \
	  -w /go/src/$(PKG) golang:1.11.5 \
	  make $(TARGET)

image-dev: target
	cd dist && cp ../Dockerfile ./ && \
	docker build -t $(TARGET_IMAGE_DEV) .

push-dev:
	docker push $(TARGET_IMAGE_DEV)

image-prd: target
	cd dist && cp ../Dockerfile ./ && \
	docker build -t $(TARGET_IMAGE_PRD) .

push-prd:
	docker push $(TARGET_IMAGE_PRD)
clean:
	rm -rf dist

.PHONY: image target clean push $(TARGET)
