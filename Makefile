BINARY=gocmd

ifneq ($(RELEASE),)
	export CGO_CPPFLAGS=$(CPPFLAGS)
	export CGO_CFLAGS=$(CFLAGS)
	export CGO_CXXFLAGS=$(CXXFLAGS)
	export CGO_LDFLAGS=$(LDFLAGS)
	FLAGS=-v -buildmode=pie -trimpath -mod=readonly -modcacherw -ldflags='-linkmode=external -w -X=github.com/tbuen/gocmd/internal/app.version=$(RELEASE) -extldflags=$(LDFLAGS)'
else
	FLAGS=-v
endif

gocmd:
	go build -o $(BINARY) $(FLAGS) ./cmd/...

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: clean
clean:
	rm -f $(BINARY)

.PHONY: install
install:
	mkdir -p $(DESTDIR)/usr/bin
	cp $(BINARY) $(DESTDIR)/usr/bin/
