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

func main() {
	imgPath := "C:/Users/pierr/Desktop/Programmation/GO/test/house.jpg"
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
	for x := 0; x <= size.X; x++ {
		for y := 0; y <= size.Y; y++ {
			pixel := img.At(x, y)
			originalColor := color.RGBAModel.Convert(pixel).(color.RGBA)
			I := float64((originalColor.R + originalColor.G + originalColor.B) / (3))
			var IH2 float64
			var IV2 float64
			//filtre horizontal
			if y > 1 {
				pixel2 := img.At(x, y-1)
				originalColor2 := color.RGBAModel.Convert(pixel2).(color.RGBA)
				IH2 = float64((originalColor2.R + originalColor2.G + originalColor2.B) / 3)
			} else {
				IH2 = I
			}
			IH := I - IH2
			if x > 1 {
				pixel2 := img.At(x-1, y)
				originalColor2 := color.RGBAModel.Convert(pixel2).(color.RGBA)
				IV2 = float64((originalColor2.R + originalColor2.G + originalColor2.B) / 3)
			} else {
				IV2 = I
			}
			IV := I - IV2
			I_tot := math.Sqrt(math.Pow(IH, 2) + math.Pow(IV, 2))
			I_final := uint8(I_tot)
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
