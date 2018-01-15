APP = k8s-publisher
GIT_COMMIT := $(shell git rev-parse HEAD)
VERSION := $(shell git describe --tags --abbrev=0 2>/dev/null)
LD_RELEASE_FLAGS += -X github.com/zmalik/k8s-publisher/pkg/version.BuildCommit=${GIT_COMMIT}
LD_RELEASE_FLAGS += -X github.com/zmalik/k8s-publisher/pkg/version.Version=${VERSION}
FOLDERS = ./pkg/...


default: test build

dep:
	dep ensure -v

build: dep test
	go build -ldflags "$(LD_RELEASE_FLAGS)" -o $(APP)

docker:
	docker build -t zmalikshxil/k8s-publisher:${VERSION} .

push: docker
	docker push zmalikshxil/k8s-publisher:${VERSION}

clean:
	rm $(APP)

test:
	go test -v  $(FOLDERS)

.PHONY: default build dockerbinary docker clean test dep