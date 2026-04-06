package infrastructure

import (
	"fractal-flame/internal/domain"
	"image"
	"image/png"
	"os"
)

func ToImage(pixels [][]domain.Pixel) *image.RGBA {
	height := len(pixels)
	width := len(pixels[0])
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			i := img.PixOffset(x, y)
			r := uint8(pixels[y][x].Color.R * 255.0)
			g := uint8(pixels[y][x].Color.G * 255.0)
			b := uint8(pixels[y][x].Color.B * 255.0)
			img.Pix[i+0] = r
			img.Pix[i+1] = g
			img.Pix[i+2] = b
			img.Pix[i+3] = 255
		}
	}
	return img
}

func SaveImage(img *image.RGBA, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)
	return png.Encode(f, img)
}
