package lsb

import (
	"encoding/binary"
	"fmt"
)

func Encode(message, container []byte, nBits int) error {
	switch nBits {
	case 1:
		return encode1bit(message, container)
	case 2:
		return encode2bit(message, container)
	case 3:
		return encode3bit(message, container)
	default:
		return fmt.Errorf("number of bits must be 1, 2 or 3")
	}
}

func Decode(container []byte, nBits int) ([]byte, error) {
	switch nBits {
	case 1:
		return decode1bit(container), nil
	case 2:
		return decode2bit(container), nil
	case 3:
		return decode3bit(container), nil
	default:
		return nil, fmt.Errorf("number of bits must be 1, 2 or 3")
	}
}

func encode1bit(message, container []byte) error {
	if (len(message)+4)*8 > len(container) {
		return fmt.Errorf("input file is too big")
	}

	dataSize := make([]byte, 4)
	binary.LittleEndian.PutUint32(dataSize, uint32(len(message)))

	message = append(dataSize, message...)
	for i := range message {
		for j := 0; j < 8; j++ {
			container[i*8+j] = (container[i*8+j] & 0b11111110) | ((message[i] >> j) & 1)
		}
	}

	return nil
}

func decode1bit(container []byte) []byte {
	dataSize := make([]byte, 4)
	for i := range dataSize {
		for j := 0; j < 8; j++ {
			dataSize[i] = dataSize[i] | ((container[i*8+j] & 1) << j)
		}
	}

	message := make([]byte, binary.LittleEndian.Uint32(dataSize))
	for i := range message {
		for j := 0; j < 8; j++ {
			message[i] = message[i] | ((container[(i+4)*8+j] & 1) << j)
		}
	}

	return message
}

func encode2bit(message, container []byte) error {
	if (len(message)+4)*8 > len(container)*2 {
		return fmt.Errorf("input file is too big")
	}

	dataSize := make([]byte, 4)
	binary.LittleEndian.PutUint32(dataSize, uint32(len(message)))

	message = append(dataSize, message...)

	for i := range message {
		container[i*4+0] = (container[i*4+0] & 0b11111100) | ((message[i] >> 0) & 0b11)
		container[i*4+1] = (container[i*4+1] & 0b11111100) | ((message[i] >> 2) & 0b11)
		container[i*4+2] = (container[i*4+2] & 0b11111100) | ((message[i] >> 4) & 0b11)
		container[i*4+3] = (container[i*4+3] & 0b11111100) | ((message[i] >> 6) & 0b11)
	}

	return nil
}

func decode2bit(container []byte) []byte {
	dataSize := make([]byte, 4)
	for i := range dataSize {
		dataSize[i] = dataSize[i] | ((container[i*4+0] & 0b11) << 0)
		dataSize[i] = dataSize[i] | ((container[i*4+1] & 0b11) << 2)
		dataSize[i] = dataSize[i] | ((container[i*4+2] & 0b11) << 4)
		dataSize[i] = dataSize[i] | ((container[i*4+3] & 0b11) << 6)

	}

	message := make([]byte, binary.LittleEndian.Uint32(dataSize))
	for i := range message {
		message[i] = message[i] | ((container[(i+4)*4+0] & 0b11) << 0)
		message[i] = message[i] | ((container[(i+4)*4+1] & 0b11) << 2)
		message[i] = message[i] | ((container[(i+4)*4+2] & 0b11) << 4)
		message[i] = message[i] | ((container[(i+4)*4+3] & 0b11) << 6)
	}

	return message
}

func encode3bit(message, container []byte) error {
	if (len(message)+4)*8 > len(container)*3 {
		return fmt.Errorf("input file is too big")
	}

	dataSize := make([]byte, 4)
	binary.LittleEndian.PutUint32(dataSize, uint32(len(message)))

	message = append(dataSize, message...)
	for i := 0; i < len(message)/3; i++ {
		container[i*8+0] = container[i*8+0]&0b11111000 | (message[i*3+0] & 0b111)
		container[i*8+1] = container[i*8+1]&0b11111000 | (message[i*3+0] >> 3 & 0b111)
		container[i*8+2] = container[i*8+2]&0b11111100 | (message[i*3+0] >> 6)
		container[i*8+2] = container[i*8+2]&0b11111011 | (message[i*3+1] << 2 & 0b100)
		container[i*8+3] = container[i*8+3]&0b11111000 | (message[i*3+1] >> 1 & 0b111)
		container[i*8+4] = container[i*8+4]&0b11111000 | (message[i*3+1] >> 4 & 0b111)
		container[i*8+5] = container[i*8+5]&0b11111110 | (message[i*3+1] >> 7)
		container[i*8+5] = container[i*8+5]&0b11111001 | (message[i*3+2] << 1 & 0b110)
		container[i*8+6] = container[i*8+6]&0b11111000 | (message[i*3+2] >> 2 & 0b111)
		container[i*8+7] = container[i*8+7]&0b11111000 | (message[i*3+2] >> 5)
	}

	switch len(message) % 3 {
	case 1:
		i := len(message) / 3
		container[i*8+0] = container[i*8+0]&0b11111000 | (message[i*3+0] & 0b111)
		container[i*8+1] = container[i*8+1]&0b11111000 | (message[i*3+0] >> 3 & 0b111)
		container[i*8+2] = container[i*8+2]&0b11111100 | (message[i*3+0] >> 6)
	case 2:
		i := len(message) / 3
		container[i*8+0] = container[i*8+0]&0b11111000 | (message[i*3+0] & 0b111)
		container[i*8+1] = container[i*8+1]&0b11111000 | (message[i*3+0] >> 3 & 0b111)
		container[i*8+2] = container[i*8+2]&0b11111100 | (message[i*3+0] >> 6)
		container[i*8+2] = container[i*8+2]&0b11111011 | (message[i*3+1] << 2 & 0b100)
		container[i*8+3] = container[i*8+3]&0b11111000 | (message[i*3+1] >> 1 & 0b111)
		container[i*8+4] = container[i*8+4]&0b11111000 | (message[i*3+1] >> 4 & 0b111)
		container[i*8+5] = container[i*8+5]&0b11111110 | (message[i*3+1] >> 7)
	}

	return nil
}

func decode3bit(container []byte) []byte {
	dataSize := make([]byte, 4)
	dataSize[0] = dataSize[0] | (container[0] & 0b111)
	dataSize[0] = dataSize[0] | (container[1] << 3 & 0b111000)
	dataSize[0] = dataSize[0] | (container[2] << 6)
	dataSize[1] = dataSize[1] | (container[2] >> 2 & 1)
	dataSize[1] = dataSize[1] | (container[3] << 1 & 0b1110)
	dataSize[1] = dataSize[1] | (container[4] << 4 & 0b1110000)
	dataSize[1] = dataSize[1] | (container[5] << 7)
	dataSize[2] = dataSize[2] | (container[5] >> 1 & 0b11)
	dataSize[2] = dataSize[2] | (container[6] << 2 & 0b11100)
	dataSize[2] = dataSize[2] | (container[7] << 5)
	dataSize[3] = dataSize[3] | (container[8] & 0b111)
	dataSize[3] = dataSize[3] | (container[9] << 3 & 0b111000)
	dataSize[3] = dataSize[3] | (container[10] << 6)

	data := make([]byte, binary.LittleEndian.Uint32(dataSize)+4)
	for i := 0; i < len(data)/3; i++ {
		data[i*3+0] = data[i*3+0] | (container[i*8+0] & 0b111)
		data[i*3+0] = data[i*3+0] | (container[i*8+1] << 3 & 0b111000)
		data[i*3+0] = data[i*3+0] | (container[i*8+2] << 6)
		data[i*3+1] = data[i*3+1] | (container[i*8+2] >> 2 & 1)
		data[i*3+1] = data[i*3+1] | (container[i*8+3] << 1 & 0b1110)
		data[i*3+1] = data[i*3+1] | (container[i*8+4] << 4 & 0b1110000)
		data[i*3+1] = data[i*3+1] | (container[i*8+5] << 7)
		data[i*3+2] = data[i*3+2] | (container[i*8+5] >> 1 & 0b11)
		data[i*3+2] = data[i*3+2] | (container[i*8+6] << 2 & 0b11100)
		data[i*3+2] = data[i*3+2] | (container[i*8+7] << 5)
	}

	switch len(data) % 3 {
	case 1:
		i := len(data) / 3
		data[i*3+0] = data[i*3+0] | (container[i*8+0] & 0b111)
		data[i*3+0] = data[i*3+0] | (container[i*8+1] << 3 & 0b111000)
		data[i*3+0] = data[i*3+0] | (container[i*8+2] << 6)
	case 2:
		i := len(data) / 3
		data[i*3+0] = data[i*3+0] | (container[i*8+0] & 0b111)
		data[i*3+0] = data[i*3+0] | (container[i*8+1] << 3 & 0b111000)
		data[i*3+0] = data[i*3+0] | (container[i*8+2] << 6)
		data[i*3+1] = data[i*3+1] | (container[i*8+2] >> 2 & 1)
		data[i*3+1] = data[i*3+1] | (container[i*8+3] << 1 & 0b1110)
		data[i*3+1] = data[i*3+1] | (container[i*8+4] << 4 & 0b1110000)
		data[i*3+1] = data[i*3+1] | (container[i*8+5] << 7)
	}

	return data[4:]
}
