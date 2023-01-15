package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
)

func main() {
	catFile, err := os.Open("C:/Users/pierr/Desktop/Programmation/GO/test/chat.png")
	if err != nil {
		log.Fatal(err)
	}
	//defer statement done after loop
	defer catFile.Close()

	imData, imType, err := image.Decode(catFile)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(imData)
	fmt.Println(imType)
	fmt.Println(imData.At(0, 0))

	cat, err := png.Decode(catFile)
	if err != nil {
		log.Fatal(err)
	}

	a, b, c, d := cat.At(0, 0).RGBA()
	fmt.Println(a)
	fmt.Println(b)
	fmt.Println(c)
	fmt.Println(d)

}
