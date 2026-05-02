package imaging

import (
	"image"
	"image/color"

	"gonum.org/v1/gonum"
	"gonum.org/v1/gonum/mat"
)

func Img2tensor(img image.Image) [][]color.Color {
	size := img.Bounds().Size()

	tensor := make([][]color.Color, size.X)
	for x := 0; x < size.X; x++ {
		for y := 0; y < size.Y; y++ {
			tensor[x] = append(tensor[x], img.At(x, y))
		}
	}

	return tensor
}

func ProcessImage(img image.Image) image.Image {
	tensor := Img2tensor(img)
	operations := []func([][]color.Color) [][]color.Color{
		grayscale,
		gaussianBlur,
	}

	for _, operation := range operations {
		tensor = operation(tensor)
	}

	return Tensor2img(tensor)
}

func clamp(v float64) float64 {
	if v < 0 {
		return 0
	} else if v > 255 {
		return 255
	}

	return v
}

func grayscale(tensor [][]color.Color) [][]color.Color {
	width := len(tensor)
	height := len(tensor[0])
	newTensor := make([][]color.Color, width)
	for i := range newTensor {
		newTensor[i] = make([]color.Color, height)
	}
	for x := range width {
		for y := range height {
			pixel := tensor[x][y]
			if pixel == nil {
				continue
			}
			r, g, b, a := tensor[x][y].RGBA()

			rf := float64(r >> 8)
			gf := float64(g >> 8)
			bf := float64(b >> 8)
			gray := uint8(rf*0.21 + gf*0.72 + bf*0.07)
			newPixel := color.RGBA{
				gray,
				gray,
				gray,
				uint8(a >> 8),
			}
			newTensor[x][y] = newPixel
		}
	}

	return newTensor
}

func spatialFilter(tensor [][]color.Color, kernel mat.Dense) [][]color.Color {
	width := len(tensor)
	height := len(tensor[0])
	newTensor := make([][]color.Color, width)
	for i := range newTensor {
		newTensor[i] = make([]color.Color, height)
	}

	kRows, kCols := kernel.Dims()
	offsetH := kRows / 2
	offsetW := kCols / 2
	for x := offsetH; x < width-offsetH; x++ {
		for y := offsetW; y < height-offsetW; y++ {
			var rSum, gSum, bSum, aSum float64
			for ka := range kRows {
				for kb := range kCols {
					// Вычисляем координаты соседа
					ix := x + ka - offsetH
					iy := y + kb - offsetW

					if tensor[ix][iy] == nil {
						continue
					}

					// Получаем цвет (RGBA() возвращает 0-65535)
					r, g, b, a := tensor[ix][iy].RGBA()
					weight := kernel.At(ka, kb)

					// Переводим в 0-255 и умножаем на вес ядра
					rSum += float64(r>>8) * weight
					gSum += float64(g>>8) * weight
					bSum += float64(b>>8) * weight
					aSum += float64(a>>8) * weight
				}
			}

			// clamping values
			newTensor[x][y] = color.RGBA{
				R: uint8(clamp(rSum)),
				G: uint8(clamp(gSum)),
				B: uint8(clamp(bSum)),
				A: uint8(clamp(aSum)),
			}
		}
	}
	return newTensor
}

func gaussianBlur(tensor [][]color.Color) [][]color.Color {
	gonum.Version()
	gaussianKernel := mat.NewDense(5, 5, []float64{
		1.0 / 256, 4.0 / 256, 6.0 / 256, 4.0 / 256, 1.0 / 256,
		4.0 / 256, 16.0 / 256, 24.0 / 256, 16.0 / 256, 4.0 / 256,
		6.0 / 256, 24.0 / 256, 36.0 / 256, 24.0 / 256, 6.0 / 256,
		4.0 / 256, 16.0 / 256, 24.0 / 256, 16.0 / 256, 4.0 / 256,
		1.0 / 256, 4.0 / 256, 6.0 / 256, 4.0 / 256, 1.0 / 256,
	})

	return spatialFilter(tensor, *gaussianKernel)
}

func Tensor2img(tensor [][]color.Color) image.Image {
	width := len(tensor)
	height := len(tensor[0])
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	size := img.Bounds().Size()
	for x := 0; x < size.X; x++ {
		for y := 0; y < size.Y; y++ {
			p := tensor[x][y]
			if p == nil {
				continue
			}

			idx := img.PixOffset(x, y)
			if c, ok := p.(color.RGBA); ok {
				img.Pix[idx] = c.R
				img.Pix[idx+1] = c.G
				img.Pix[idx+2] = c.B
				img.Pix[idx+3] = c.A
			} else {
				// converting if pixel is not RGBA
				r, g, b, a := p.RGBA()
				img.Pix[idx] = uint8(r >> 8)
				img.Pix[idx+1] = uint8(g >> 8)
				img.Pix[idx+2] = uint8(b >> 8)
				img.Pix[idx+3] = uint8(a >> 8)
			}

		}
	}
	return img
}
