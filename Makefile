TARGET=news_watch_notice
PKG=$(TARGET)
TAG=latest
IMAGE_PREFIX?=hanyajun
IMAGE_PREFIX_PRD=hanyajun
TARGET_IMAGE_DEV=$(IMAGE_PREFIX)/$(TARGET):$(TAG)
TARGET_IMAGE_PRD=$(IMAGE_PREFIX_PRD)/$(TARGET):$(TAG)
image-dev:
	docker build -t $(TARGET_IMAGE_DEV) .
build-dev-nch:
	docker build -t $(TARGET_IMAGE_DEV) --no-cache .
push-dev:
	docker push $(TARGET_IMAGE_DEV)

image-prd:
	docker build -t $(TARGET_IMAGE_PRD) .

push-prd:
	docker push $(TARGET_IMAGE_PRD)


ci-dev: image-dev push-dev
ci-dev-nch: build-dev-nch push-dev
clean:
	rm -rf dist

.PHONY: image target clean push $(TARGET)
