all: build encode decode

encode:
	./lsb encode 3 res/image.png res/girl.bmp out/encoded.bmp

decode:
	./lsb decode 3 out/encoded.bmp out/decoded

build:
	@go build -o lsb .

run:
	@./lsb
