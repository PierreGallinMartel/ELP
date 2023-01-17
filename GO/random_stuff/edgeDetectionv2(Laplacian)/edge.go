package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"math"
	"os"
	"path/filepath"
	"strings"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func Convolution(img *image.Image, matrice [][]int) *image.NRGBA {
	imageRGBA := image.NewNRGBA((*img).Bounds())
	w := (*img).Bounds().Dx()
	h := (*img).Bounds().Dy()
	sumR := 0
	sumB := 0
	sumG := 0
	var r uint32
	var g uint32
	var b uint32
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {

			for i := -1; i <= 1; i++ {
				for j := -1; j <= 1; j++ {

					var imageX int
					var imageY int

					imageX = x + i
					imageY = y + j

					r, g, b, _ = (*img).At(imageX, imageY).RGBA()
					sumG = (sumG + (int(g) * matrice[i+1][j+1]))
					sumR = (sumR + (int(r) * matrice[i+1][j+1]))
					sumB = (sumB + (int(b) * matrice[i+1][j+1]))
				}
			}

			imageRGBA.Set(x, y, color.NRGBA{
				uint8(min(sumR/9, 0xffff) >> 8),
				uint8(min(sumG/9, 0xffff) >> 8),
				uint8(min(sumB/9, 0xffff) >> 8),
				255,
			})

			sumR = 0
			sumB = 0
			sumG = 0

		}
	}

	return imageRGBA
}

func main() {
	imgPath := "C:/Users/pierr/Desktop/Programmation/ELP/ELP/GO/random_stuff/images/cathedral_blurred.jpg"
	f, err := os.Open(imgPath)
	check(err)
	defer f.Close()
	img, imType, imerr := image.Decode(f)
	print(imType)
	print(imerr)
	size := img.Bounds().Size()
	rect := image.Rect(0, 0, size.X, size.Y)
	wImg := image.NewRGBA(rect)

	//filtrage
	var s float64
	s = 80
	distance := 1
	for x := 0; x <= size.X; x++ {
		for y := 0; y <= size.Y; y++ {
			pixel := img.At(x, y)
			originalColor := color.RGBAModel.Convert(pixel).(color.RGBA)
			C := float64((originalColor.R + originalColor.G + originalColor.B) / (3))
			var N float64
			var S float64
			var E float64
			var W float64
			var laplacien float64
			var I_final uint8
			//filtre vertical
			if y > distance {
				pixel2 := img.At(x, y-distance)
				originalColor2 := color.RGBAModel.Convert(pixel2).(color.RGBA)
				N = float64((originalColor2.R + originalColor2.G + originalColor2.B) / 3)
			} else {
				N = C
			}
			if y < size.Y-distance {
				pixel2 := img.At(x, y+distance)
				originalColor2 := color.RGBAModel.Convert(pixel2).(color.RGBA)
				S = float64((originalColor2.R + originalColor2.G + originalColor2.B) / 3)
			} else {
				S = C
			}
			CV := N + S - 2*C
			if x > distance {
				pixel2 := img.At(x-distance, y)
				originalColor2 := color.RGBAModel.Convert(pixel2).(color.RGBA)
				W = float64((originalColor2.R + originalColor2.G + originalColor2.B) / 3)
			} else {
				W = C
			}
			if x < size.X-distance {
				pixel2 := img.At(x+distance, y)
				originalColor2 := color.RGBAModel.Convert(pixel2).(color.RGBA)
				E = float64((originalColor2.R + originalColor2.G + originalColor2.B) / 3)
			} else {
				E = C
			}
			CH := W + E - 2*C
			laplacien = (CV + CH)
			if math.Abs(laplacien) < s {
				I_final = 0
			} else {
				I_final = 255
			}
			c := color.RGBA{
				R: I_final, G: I_final, B: I_final, A: 1,
			}
			wImg.Set(x, y, c)
		}

	}
	ext := filepath.Ext(imgPath)
	name := strings.TrimSuffix(filepath.Base(imgPath), ext)
	newImagePath := fmt.Sprintf("%s/%s_edged%s", filepath.Dir(imgPath), name, ext)
	fg, err := os.Create(newImagePath)
	defer fg.Close()
	check(err)
	err = jpeg.Encode(fg, wImg, nil)
	check(err)

}
