VERSION  ?= 1.0.0
ARCH     := amd64
BINARY   := clarocr
DEB_DIR  := dist/$(BINARY)_$(VERSION)_$(ARCH)
DEB_FILE := dist/$(BINARY)_$(VERSION)_$(ARCH).deb

.PHONY: build deb deps clean

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

clean:
	rm -rf $(BINARY) dist/
