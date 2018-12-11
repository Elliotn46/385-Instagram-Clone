

VERSION=latest
PROJECT=$(shell gcloud config get-value project)
REGISTRY=gcr.io/$(PROJECT)

.PHONY: accountauth accountsearch rabbitconsumer pushimages all

accountauth:

        docker build -t account_auth:$(VERSION) -f accounts/Dockerfile accounts/
        docker tag account_auth:$(VERSION) $(REGISTRY)/account_auth:$(VERSION)

rabbitconsumer:

        docker build -t rabbit_consumer:$(VERSION) -f consumer/rabbit/Dockerfile consumer/rabbit
        docker tag rabbit_consumer:$(VERSION) $(REGISTRY)/rabbit_consumer:$(VERSION)

accountsearch:

        docker build -t account_search:$(VERSION) -f search/Dockerfile search/
        docker tag account_search:$(VERSION) $(REGISTRY)/account_search:$(VERSION)

pushimages: accountsearch rabbitconsumer accountauth

        gcloud docker -- push $(REGISTRY)/account_auth:$(VERSION)
        gcloud docker -- push $(REGISTRY)/account_search:$(VERSION)
        gcloud docker -- push $(REGISTRY)/rabbit_consumer:$(VERSION)

all: pushimages
