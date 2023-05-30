package bmp

import (
	"encoding/binary"
	"encoding/json"
	"log"
	"os"
)

const (
	RedOffset   = 2
	GreenOffset = 1
	BlueOffset  = 0
)

type BitmapFileHeader struct {
	FileType  uint16
	Size      uint32
	Reserved1 uint16
	Reserved2 uint16
	Offbits   uint32
}

type BitmapInfoHeader struct {
	Size          uint32
	Width         int32
	Height        int32
	Planes        uint16
	BitCount      uint16
	Compression   uint32
	SizeImage     uint32
	XPelsPerMeter int32
	YPelsPerMeter int32
	ClrUsed       uint32
	ClrImportant  uint32
}

type Image struct {
	FileHeader BitmapFileHeader
	InfoHeader BitmapInfoHeader
	PixelData  []byte
}

func NewImage(filename string) (*Image, error) {
	f, err := os.Open(filename)
	if err != nil {
        return nil, err
	}

	img := new(Image)

	binary.Read(f, binary.LittleEndian, &img.FileHeader)
	binary.Read(f, binary.LittleEndian, &img.InfoHeader)

	bytesPerLine := img.InfoHeader.Width * 3
	nDummyBytes := bytesPerLine % 4

	nPixelBytes := bytesPerLine * Abs(img.InfoHeader.Height)
	img.PixelData = make([]byte, nPixelBytes)

	f.Seek(int64(img.FileHeader.Offbits), 0)

	for i := int32(0); i < Abs(img.InfoHeader.Height); i++ {
		curInd := i * bytesPerLine
		f.Read(img.PixelData[curInd : curInd+bytesPerLine])
		f.Seek(int64(nDummyBytes), 1)
	}

	f.Close()

	return img, nil
}

func (bi *Image) Width() int {
	return int(bi.InfoHeader.Width)
}

func (bi *Image) Height() int {
	return int(bi.InfoHeader.Height)
}

func (bi *Image) WriteToFile(filename string) {
	f, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}

	binary.Write(f, binary.LittleEndian, &bi.FileHeader)
	binary.Write(f, binary.LittleEndian, &bi.InfoHeader)

	f.Seek(int64(bi.FileHeader.Offbits), 0)

	bytesPerLine := bi.InfoHeader.Width
	if bi.InfoHeader.BitCount == 24 {
		bytesPerLine *= 3
	}
	dummyBytes := make([]byte, bytesPerLine%4)

	for i := int32(0); i < Abs(bi.InfoHeader.Height); i++ {
		curInd := i * bytesPerLine
		f.Write(bi.PixelData[curInd : curInd+bytesPerLine])
		f.Write(dummyBytes)
	}

	f.Close()
}

func (bi *Image) String() string {
	b, _ := json.MarshalIndent(bi.FileHeader, "", "  ")
	kek1 := string(b)
	b, _ = json.MarshalIndent(bi.InfoHeader, "", "  ")
	kek2 := string(b)
	return kek1 + kek2
}

func (bi *Image) getComponentBytes(off int) []byte {
	n := len(bi.PixelData) / 3
	data := make([]byte, n)
	for i := 0; i < n; i++ {
		data[i] = bi.PixelData[i*3+off]
	}
	return data
}

func (bi *Image) GetRgb() ([]byte, []byte, []byte) {
	return bi.GetRedBytes(), bi.GetGreenBytes(), bi.GetBlueBytes()
}

func (bi *Image) GetYCbCr() ([]byte, []byte, []byte) {
	return RgbToYcbcr(bi.GetRgb())
}

func (bi *Image) GetRedBytes() []byte {
	return bi.getComponentBytes(RedOffset)
}

func (bi *Image) GetGreenBytes() []byte {
	return bi.getComponentBytes(GreenOffset)
}

func (bi *Image) GetBlueBytes() []byte {
	return bi.getComponentBytes(BlueOffset)
}

func (bi *Image) SetRgb(red, green, blue []byte) {
	for i := 0; i < len(red); i++ {
		ind := i * 3
		bi.PixelData[ind+RedOffset] = red[i]
		bi.PixelData[ind+GreenOffset] = green[i]
		bi.PixelData[ind+BlueOffset] = blue[i]
	}
}

func (bi *Image) SetYCbCr(Y, Cb, Cr []byte) {
	bi.SetRgb(YcbcrToRgb(Y, Cb, Cr))
}
