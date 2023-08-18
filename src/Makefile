 # Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
BINARY_NAME=../bin/swag2nginx
PREFIX=/usr/local

all: build

build: 
	$(GOBUILD) -o $(BINARY_NAME) genconfig/main.go

clean:
	rm -f $(BINARY_NAME)
	rm -f $(PREFIX)/bin/swag2nginx

install:
	install $(BINARY_NAME) $(PREFIX)/bin/swag2nginx 
