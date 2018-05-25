OS = darwin freebsd linux openbsd windows
ARCHS = 386 arm amd64 arm64

all: build release

build: deps
	go build

release: clean deps
	@for arch in $(ARCHS);\
	do \
		for os in $(OS);\
		do \
			echo "Building $$os-$$arch"; \
			mkdir -p build/webhook-$$os-$$arch/; \
			GOOS=$$os GOARCH=$$arch CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static"' -o build/webhook-$$os-$$arch/webhook; \
			tar cz -C build -f build/webhook-$$os-$$arch.tar.gz webhook-$$os-$$arch; \
		done \
	done

test: deps
	mkdir -p .cover
	go test -covermode=count -cover -coverprofile .cover/cover.out ./...

deps:
	go get -d -v -t ./...

clean:
	rm -rf build
	rm -f webhook
