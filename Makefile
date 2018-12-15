VERSION ?= latest
PROJECT ?= $(shell gcloud config get-value project)
REGISTRY ?= gcr.io/$(PROJECT)

images := accounts/auth accounts/search content/writer consumer/rabbit
tags := $(patsubst %,$(REGISTRY)/%:$(VERSION),$(images))

.PHONY: buildimages $(images) pushimages

buildimages: $(images)
	@echo 'Images built, now `make pushimages` to push.'

$(images):
	docker build -t $@:$(VERSION) -f $@/Dockerfile $@
	docker tag $@:$(VERSION) $(REGISTRY)/$@:$(VERSION)

pushimages: buildimages
	for tag in $(tags); do gcloud docker -- push $$tag; done
