<mkconfig

all: build

build :V:
	go build -o exe/ ./cmd/...

clean :V:
	rm -f exe/*

