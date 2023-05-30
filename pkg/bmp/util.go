package bmp

import (
	"fmt"
	"io"
	"math"

	"golang.org/x/exp/constraints"
)

func Abs[T constraints.Signed](x T) T {
	if x > 0 {
		return x
	}
	return -x
}

func Mean(data []byte, from, to int) float64 {
	var sum float64
	for i := from; i < to; i++ {
		sum += float64(data[i])
	}
	return sum / float64(to-from)
}

func StdDeviation(data []byte, from, to int) float64 {
	m := Mean(data, from, to)
	var sum float64
	for i := from; i < to; i++ {
		sum += math.Pow(float64(data[i])-m, 2)
	}
	return math.Sqrt(1.0 / (float64(to-from) - 1) * sum)
}

func CorrelationCoefficient(a, b []byte, from, to int) float64 {
	mA := Mean(a, from, to)
	mB := Mean(b, from, to)
	sdA := StdDeviation(a, from, to)
	sdB := StdDeviation(b, from, to)
	var sum float64
	for i := from; i < to; i++ {
		sum += (float64(a[i]) - mA) * (float64(b[i]) - mB)
	}

	return sum / (float64(to-from) * sdA * sdB)
}

func sample1(data []byte, x, y, width, height int) []byte {
	result := make([]byte, len(data))

	ind := func(i, j int) int {
		return i*width + j
	}

	if y >= 0 {
		if x >= 0 {
			for j := 0; j < width-x; j++ {
				for i := 0; i < height-y; i++ {
					result[ind(i, j)] = data[ind(i, j)]
				}
			}
		} else {
			for j := -x; j < width; j++ {
				for i := 0; i < height-y; i++ {
					result[ind(i, j+x)] = data[ind(i, j)]
				}
			}
		}
	} else {
		if x >= 0 {
			for j := 0; j < width-x; j++ {
				for i := -y; i < height; i++ {
					result[ind(i+y, j)] = data[ind(i, j)]
				}
			}
		} else {
			for j := -x; j < width; j++ {
				for i := -y; i < height; i++ {
					result[ind(i+y, j+x)] = data[ind(i, j)]
				}
			}
		}
	}
	return result
}

func sample2(data []byte, x, y, width, height int) []byte {
	result := make([]byte, len(data))

	ind := func(i, j int) int {
		return i*width + j
	}

	if y >= 0 {
		if x >= 0 {
			for j := x; j < width; j++ {
				for i := y; i < height; i++ {
					result[ind(i-y, j-x)] = data[ind(i, j)]
				}
			}
		} else {
			for j := 0; j < width+x; j++ {
				for i := y; i < height; i++ {
					result[ind(i-y, j)] = data[ind(i, j)]
				}
			}
		}
	} else {
		if x >= 0 {
			for j := x; j < width; j++ {
				for i := 0; i < height+y; i++ {
					result[ind(i, j-x)] = data[ind(i, j)]
				}
			}
		} else {
			for j := 0; j < width+x; j++ {
				for i := 0; i < height+y; i++ {
					result[ind(i, j)] = data[ind(i, j)]
				}
			}
		}
	}
	return result
}

func Autocorrelation(data []byte, y, width, height int, out io.Writer) {
	maxX := width / 4
	minX := -maxX
	stepX := 5

	for x := minX; x < maxX; x += stepX {
		a1 := sample1(data, x, y, width, height)
		a2 := sample2(data, x, y, width, height)
		r := CorrelationCoefficient(a1, a2, 0, len(a1))
		fmt.Fprintf(out, "%d %f\n", x, r)
	}
}

func Clip[T constraints.Integer | constraints.Float](x T) byte {
	result := math.Round(float64(x))
	if result < 0 {
		return 0
	} else if result > math.MaxUint8 {
		return math.MaxUint8
	}
	return byte(x)
}

func RgbToYcbcr(R, G, B []byte) (Y, Cb, Cr []byte) {
	n := len(R)
	Y, Cb, Cr = make([]byte, n), make([]byte, n), make([]byte, n)
	for i := 0; i < len(R); i++ {
		r, g, b := float64(R[i]), float64(G[i]), float64(B[i])
		y := 0.299*r + 0.587*g + 0.114*b

		Y[i] = Clip(y)
		Cb[i] = Clip(0.5643*(b-y) + 128)
		Cr[i] = Clip(0.7132*(r-y) + 128)
	}
	return Y, Cb, Cr
}

func YcbcrToRgb(Y, Cb, Cr []byte) (R, G, B []byte) {
	n := len(Y)
	R, G, B = make([]byte, n), make([]byte, n), make([]byte, n)
	for i := 0; i < len(Y); i++ {
		y, cb, cr := float64(Y[i]), float64(Cb[i]), float64(Cr[i])

		R[i] = Clip(y + 1.402*(cr-128))
		G[i] = Clip(y - 0.714*(cr-128) - 0.334*(cb-128))
		B[i] = Clip(y + 1.772*(cb-128))
	}
	return R, G, B
}

func Psnr(a, b []byte) float64 {
	numerator := float64(len(a)) * math.Pow(float64(math.MaxUint8), 2)
	var denominator float64
	for i := 0; i < len(a); i++ {
		denominator += math.Pow(float64(a[i])-float64(b[i]), 2)
	}
	return 10.0 * math.Log10(numerator/denominator)
}

func Frequency[T constraints.Ordered](data []T) map[T]int {
	freq := make(map[T]int)
	for i := range data {
		freq[data[i]]++
	}
	return freq
}

func Entropy[T constraints.Integer](data []T) float64 {
	freq := Frequency(data)
	from, to := T(0), T(0)
	for k := range freq {
		if k < from {
			from = k
		} else if k > to {
			to = k
		}
	}

	n := float64(len(data))

	var entropy float64
	for x := from; x < to; x++ {
		pX := float64(freq[x]) / n
		if pX == 0 {
			continue
		}
		entropy += pX * math.Log2(pX)
	}

	return -entropy
}
