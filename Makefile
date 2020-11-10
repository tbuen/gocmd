BINARY=gocmd

ifneq ($(RELEASE),)
	export CGO_CPPFLAGS=$(CPPFLAGS)
	export CGO_CFLAGS=$(CFLAGS)
	export CGO_CXXFLAGS=$(CXXFLAGS)
	export CGO_LDFLAGS=$(LDFLAGS)
	FLAGS=-v -tags release -buildmode=pie -trimpath -mod=readonly -modcacherw -ldflags='-linkmode=external -w -X=github.com/tbuen/gocmd/internal/app.version=$(RELEASE) -extldflags=$(LDFLAGS)'
else
	FLAGS=-v
endif

.PHONY: gocmd
gocmd:
	go build -o $(BINARY) $(FLAGS) ./cmd/...

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: clean
clean:
	@rm -f $(BINARY)
	@rm -rf build/package/gocmd* build/package/pkg build/package/src

.PHONY: install
install:
	mkdir -p $(DESTDIR)/usr/bin
	cp $(BINARY) $(DESTDIR)/usr/bin/
	mkdir -p $(DESTDIR)/usr/share/gocmd/config
	cp configs/apps.xml $(DESTDIR)/usr/share/gocmd/config/apps.xml.sample
