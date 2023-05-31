package main

import (
	"bufio"
	"crypto/rand"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/elvekdarzhinov/lsb/pkg/bmp"
	"github.com/elvekdarzhinov/lsb/pkg/lsb"
)

const (
	helpMessage = `USAGE:
  ENCODE:
    lsb encode n_bits message_file in_img.bmp out_img.bmp

  DECODE:
    lsb decode n_bits in_img.bmp out_file

  n_bits - number of least significant bits to use (1, 2 or 3)`
)

func main() {
	defer handleErr()
	if len(os.Args) > 1 && os.Args[1] == "test" {
		test()
		return
	}
	run()
}

func test() {
	img, _ := bmp.NewImage("res/girl.bmp")

	nBits := 3

	generateRandomFile("test/data.txt", (len(img.PixelData)/8*nBits-4))
	message, _ := os.ReadFile("test/data.txt")

	lsb.Encode(message, img.PixelData, nBits)
	img.WriteToFile("out/encoded.bmp")

	tmp, _ := bmp.NewImage("res/girl.bmp")
	fmt.Println(bmp.Psnr(img.PixelData, tmp.PixelData))
}

func run() {
	if !(len(os.Args) == 5 || len(os.Args) == 6) {
		fmt.Println(helpMessage)
		return
	}

	nBits, err := strconv.Atoi(os.Args[2])
	if err != nil {
		panic(err)
	}

	switch os.Args[1] {
	case "encode":
		message, err := os.ReadFile(os.Args[3])
		if err != nil {
			panic(err)
		}

		img, err := bmp.NewImage(os.Args[4])
		if err != nil {
			panic(err)
		}

		err = lsb.Encode(message, img.PixelData, nBits)
		if err != nil {
			panic(err)
		}

		img.WriteToFile(os.Args[5])

	case "decode":
		img, err := bmp.NewImage(os.Args[3])
		if err != nil {
			panic(err)
		}

		message, err := lsb.Decode(img.PixelData, nBits)
		if err != nil {
			panic(err)
		}

		os.WriteFile(os.Args[4], message, os.ModePerm)

	default:
		panic("unknown command: " + os.Args[1])
	}
}

func handleErr() {
	if err := recover(); err != nil {
		fmt.Println(err)
	}
}

func generateRandomFile(file string, size int) {
	f, _ := os.Create(file)
	w := bufio.NewWriter(f)
	defer func() { w.Flush(); f.Close() }()

	io.CopyN(w, rand.Reader, int64(size))
}

func equalFiles(fileA, fileB string) bool {
	a, _ := os.ReadFile(fileA)
	b, _ := os.ReadFile(fileB)

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
