package main

import (
	"image/color"
	"image/jpeg"
	"os"
)

func main() {
	type Changeable interface {
		Set(x, y int, c color.Color)
	}

	imgfile, err := os.Open("C:/Users/pierr/Desktop/Programmation/GO/test/chat.jpg")
	if err != nil {
		panic(err.Error())
	}
	defer imgfile.Close()

	img, err := jpeg.Decode(imgfile)
	if err != nil {
		panic(err.Error())
	}

	if cimg, ok := img.(Changeable); ok {
		// cimg is of type Changeable, you can call its Set() method (draw on it)
		cimg.Set(0, 0, color.RGBA{85, 165, 34, 255})
		cimg.Set(0, 1, color.RGBA{255, 0, 0, 255})
		// when done, save img as usual
	} else {
		// No luck... see your options below
	}
}
