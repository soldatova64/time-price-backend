# Переменные
IMAGE_NAME = soldatova64/time-price-backend
VERSION = latest
REGISTRY = docker.io/soldatovadew

.PHONY: build push

build:
	docker build -t $(IMAGE_NAME):$(VERSION) .

push: build
	docker tag $(IMAGE_NAME):$(VERSION) $(REGISTRY)/$(IMAGE_NAME):$(VERSION)
	docker push $(REGISTRY)/$(IMAGE_NAME):$(VERSION)