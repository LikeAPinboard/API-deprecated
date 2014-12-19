PRODUCT=likeapinboard
OS = $(shell uname)

all: darwin linux windows

darwin:
	GOPATH=`pwd` GOOS=darwin GOARC=amd64 go build -o build/darwin/$(PRODUCT) ./main.go
linux:
	GOPATH=`pwd` GOOS=linux GOARC=amd64 go build -o build/linux/$(PRODUCT) ./main.go
windows:
	GOPATH=`pwd` GOOS=windows GOARC=amd64 go build -o build/windows/$(PRODUCT).exe ./main.go

clean:
	rm build/darwin/$(PRODUCT)
	rm build/linux/$(PRODUCT)
	rm build/windows/$(PRODUCT).exe

install:
ifeq ($(OS),Linux)
# for Linux
	cp build/linux/$(PRODUCT) /usr/local/bin/$(PRODUCT)
	cp config.json /etc/$(PROUCT).json
	cp scripts/$(PRODUCT) /etc/init.d/$(PRODUCT)
	chmod +x /etc/init.d/$(PRODUCT)
endif
ifeq ($(OS),Darwin)
# for MacOSX
	cp build/darwin/$(PRODUCT) /usr/local/bin/$(PRODUCT)
	cp config.json /etc/$(PRODUCT).json
endif

uninstall:
ifeq ($(OS),Linux)
# for Linux
	rm -f /usr/local/bin/$(PRODUCT)
	rm -f /etc/$(PRODUCT).json
	rm -f /etc/init.d/$(PRODUCT)
endif
ifeq ($(OS),Darwin)
# for MacOSX
	rm -f /usr/local/bin/$(PRODUCT)
	rm -f /etc/$(PRODUCT).json
endif
