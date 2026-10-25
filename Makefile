VERSION  ?= 0.0.1
RELEASE  ?= 1
ARCH     := amd64
BINARY   := clarocr
DEB_DIR  := dist/$(BINARY)_$(VERSION)_$(ARCH)
DEB_FILE := dist/$(BINARY)_$(VERSION)_$(ARCH).deb

.PHONY: build deb deps rpm rpm-deps clean

deps:
	@which go > /dev/null 2>&1 || (sudo apt-get update && sudo apt-get install -y golang-go)
	sudo apt-get install -y libtesseract-dev libleptonica-dev libayatana-appindicator3-dev

build:
	go build -o $(BINARY) .

deb: deps build
	rm -rf $(DEB_DIR)
	mkdir -p $(DEB_DIR)/DEBIAN
	mkdir -p $(DEB_DIR)/usr/local/bin
	cp $(BINARY) $(DEB_DIR)/usr/local/bin/
	cp packaging/debian/control $(DEB_DIR)/DEBIAN/control
	cp packaging/debian/postinst $(DEB_DIR)/DEBIAN/postinst
	sed -i "s/^Version:.*/Version: $(VERSION)/" $(DEB_DIR)/DEBIAN/control
	chmod 755 $(DEB_DIR)/DEBIAN/postinst
	dpkg-deb --build $(DEB_DIR) $(DEB_FILE)
	@echo ""
	@echo "Pacote gerado: $(DEB_FILE)"
	@echo "Para instalar: sudo apt install ./$(DEB_FILE)"

# Explicit setup: make rpm itself never installs packages or invokes sudo.
rpm-deps:
	sudo dnf install -y --setopt=install_weak_deps=False rpm-build golang gcc gcc-c++ make pkgconf-pkg-config tesseract-devel leptonica-devel libayatana-appindicator-gtk3-devel tar gzip findutils diffutils

rpm:
	@sh packaging/rpm/build.sh "$(VERSION)" "$(RELEASE)"

clean:
	rm -rf $(BINARY) dist/
