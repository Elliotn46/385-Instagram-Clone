VERSION=latest
PROJECT=$(shell gcloud config get-value project)
REGISTRY=gcr.io/$(PROJECT)
.PHONY: account-auth
account-auth:
        docker build -t account_auth:$(VERSION) -f accounts/Dockerfile accounts/
        docker tag account_auth:$(VERSION) $(REGISTRY)/account_auth:$(VERSION)
.PHONY: rabbit-consumer
rabbit-consumer:
        docker build -t rabbit_consumer:$(VERSION) -f consumer/rabbit/Dockerfile consumer/rabbit
        docker tag rabbit_consumer:$(VERSION) $(REGISTRY)/rabbit_consumer:$(VERSION)
.PHONY: account-search
account-search:
        docker build -t account_search:$(VERSION) -f search/Dockerfile search/
        docker tag account_search:$(VERSION) $(REGISTRY)/account_search:$(VERSION)
.PHONY: push-images
push-images:
        gcloud docker -- push $(REGISTRY)/account_auth:$(VERSION)
        gcloud docker -- push $(REGISTRY)/account_search:$(VERSION)
        gcloud docker -- push $(REGISTRY)/rabbit_consumer:$(VERSION)
