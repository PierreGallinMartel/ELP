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

func main() {
	// : All your code will go here
	imgPath := "C:/Users/pierr/Desktop/Programmation/ELP/ELP/GO/random_stuff/images/cathedral.jpg"
	f, err := os.Open(imgPath)
	check(err)
	defer f.Close()
	img, imType, imerr := image.Decode(f)
	print(imType)
	print(imerr)
	size := img.Bounds().Size()
	rect := image.Rect(0, 0, size.X, size.Y)
	wImg := image.NewRGBA(rect)
	// loop though all the x
	for x := 0; x < size.X; x++ {
		// and now loop thorough all of this x's y
		for y := 0; y < size.Y; y++ {
			pixel := img.At(x, y)
			//ri, gi, bi, a := pixel.RGBA()
			originalColor := color.RGBAModel.Convert(pixel).(color.RGBA)
			//Offset colors a little, adjust it to your taste
			r := float64(originalColor.R)
			g := float64(originalColor.G)
			b := float64(originalColor.B)
			// average
			grey := uint8((r + g + b) / 3)
			c := color.RGBA{
				R: grey, G: grey, B: grey, A: uint8(originalColor.A),
			}
			wImg.Set(x, y, c)
		}
	}
	ext := filepath.Ext(imgPath)
	name := strings.TrimSuffix(filepath.Base(imgPath), ext)
	newImagePath := fmt.Sprintf("%s/%s_gray%s", filepath.Dir(imgPath), name, ext)
	fg, err := os.Create(newImagePath)
	defer fg.Close()
	check(err)
	err = jpeg.Encode(fg, wImg, nil)
	check(err)

}
