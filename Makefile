.PHONY: all windows linux

all: windows linux

windows:
	GOOS=windows GOARCH=amd64 go build -o GoWebp-windows.exe ./cmd/GoWebp

linux:
	GOOS=linux GOARCH=amd64 go build -o GoWebp-linux ./cmd/GoWebp

clean:
	rm -f GoWebp-windows.exe GoWebp-linux
