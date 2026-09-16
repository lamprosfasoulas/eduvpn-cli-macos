BINARY := eduvpn-cli
PREFIX := /usr/local/bin

.PHONY: build install uninstall clean

build:
	go build -o $(BINARY) .

install: build
	install -m 0755 $(BINARY) $(PREFIX)/$(BINARY)

uninstall:
	rm -f $(PREFIX)/$(BINARY)

clean:
	rm -f $(BINARY)
