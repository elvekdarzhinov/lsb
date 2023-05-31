all: build encode decode

t: build
	./lsb test

encode:
	./lsb encode 2 test/data.txt res/girl.bmp out/encoded.bmp

decode:
	./lsb decode 2 out/encoded.bmp out/decoded

build:
	@go build -o lsb .

run:
	@./lsb

windows:
	GOOS=windows GOARCH=amd64 go build .

mac:
	GOOS=darwin GOARCH=amd64 go build .

