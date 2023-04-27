#TAG = test-$(shell git log -1 --format=%h)
TAG = latest
WORK_DIR = .
REGISTRY = registry.cn-shanghai.aliyuncs.com/codev

image:
	docker build --target plant-shutter -t $(REGISTRY)/plant-shutter:$(TAG) -f dockerfiles/plant-shutter.dockerfile $(WORK_DIR)

push:
	docker push $(REGISTRY)/plant-shutter:$(TAG)

