package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
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
	imgPath := "C:/Users/pierr/Desktop/Programmation/ELP/ELP/GO/random_stuff/images/cathedral.jpg"
	f, err := os.Open(imgPath)
	check(err)
	defer f.Close()
	img, imType, imerr := image.Decode(f)
	print(imType)
	print(imerr)
	//size := img.Bounds().Size()
	//rect := image.Rect(0, 0, size.X, size.Y)
	kernell := [][]int{{1, 1, 1}, {1, 1, 1}, {1, 1, 1}}
	wImg := Convolution(&img, kernell)

	ext := filepath.Ext(imgPath)
	name := strings.TrimSuffix(filepath.Base(imgPath), ext)
	newImagePath := fmt.Sprintf("%s/%s_blurred%s", filepath.Dir(imgPath), name, ext)
	fg, err := os.Create(newImagePath)
	defer fg.Close()
	check(err)
	err = jpeg.Encode(fg, wImg, nil)
	check(err)

}
