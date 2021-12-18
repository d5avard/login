all:
	make clean
	make build 
	make run
test:
	ENV=TEST WD=$(PWD) go test -v ./...

cover:
	ENV=TEST WD=$(PWD) go test -v ./... -coverprofile=cover.txt

build:
	go build -o $(GOBUILDOUTPUT) -ldflags "-X main.Version=v$(VERSION)"

run:
	./$(GOBUILDOUTPUT)

clean:
	rm -f $(GOBUILDOUTPUT)

ds:
	minikube start --driver=virtualbox --container-runtime=docker
	eval $(minikube docker-env)

dr:
	docker run --rm -d -p 127.0.0.1:8080:8080 northamerica-northeast1-docker.pkg.dev/danysavard-ca-327813/diary/brbff:$(VERSION)

# k8 create
kc: kbc kbd

# k8 all
ka: kbc kpc

# k8 build container
kbc:
	docker build . \
		--build-arg VERSION=$(VERSION) \
		--build-arg GOBUILDOUTPUT=$(GOBUILDOUTPUT) \
		--rm --tag northamerica-northeast1-docker.pkg.dev/danysavard-ca-327813/diary/brbff:$(VERSION)

# k8 docker push
kpc: 
	docker push northamerica-northeast1-docker.pkg.dev/danysavard-ca-327813/diary/brbff:$(VERSION)

# k8 push deployment
kbd:
	kubectl create deployment brbff --image=northamerica-northeast1-docker.pkg.dev/danysavard-ca-327813/diary/brbff:$(VERSION)

# k8 update deployment 
kup:
	kubectl set image deployment/brbff brbff=northamerica-northeast1-docker.pkg.dev/danysavard-ca-327813/diary/brbff:$(VERSION)