TARGET = ./bin/go_http_xxx
INSTALL_DIR = /home/server/go_http_xxx/bin/

BuildVersion = v2
BuildTime = $(shell date +"%Y-%m-%d %H:%M:%S")
BuildSum = $(shell git rev-parse --short HEAD)

GO_SOURCE = $(wildcard ./src/cmd/*.go)

all: $(TARGET)
build: $(TARGET)

$(TARGET): $(GO_SOURCE)
	mkdir -p bin
	go build -ldflags "-X 'main.BuildVersion=$(BuildVersion)' -X 'main.BuildTime=$(BuildTime)' -X 'main.BuildSum=$(BuildSum)'" -o $@ $^ 

clean:
	rm -rf $(TARGET)

install:
	mkdir -p $(INSTALL_DIR)
	cp $(TARGET) $(INSTALL_DIR)