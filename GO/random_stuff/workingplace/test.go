package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var wg sync.WaitGroup

func worker(img *image.Image, wImg *image.RGBA, jobs <-chan [2]int) {
	defer wg.Done()
	cords := <-jobs
	pixel := (*img).At(cords[0], cords[1])

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
	wImg.Set(cords[0], cords[1], c)

}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	const numJobs = 1000
	jobs := make(chan [2]int, numJobs)
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
	//generate our workers
	/*
		for x := 0; x < size.X; x++ {
			for y := 0; y < size.Y; y++ {
				go worker(x, y, wImg, jobs, results)
			}
		}
	*/
	for w := 0; w <= size.Y; w++ {
		wg.Add(1)
		go worker(&img, wImg, jobs)
	}
	// loop though all the x
	go func() {
		for x := 0; x < size.X; x++ {
			// and now loop thorough all of this x's y
			for y := 0; y < size.Y; y++ {
				cords := [2]int{x, y}
				jobs <- cords
			}
		}
	}()
	wg.Wait()
	ext := filepath.Ext(imgPath)
	name := strings.TrimSuffix(filepath.Base(imgPath), ext)
	newImagePath := fmt.Sprintf("%s/%s_gray%s", filepath.Dir(imgPath), name, ext)
	fg, err := os.Create(newImagePath)
	defer fg.Close()
	check(err)
	err = jpeg.Encode(fg, wImg, nil)
	check(err)

}
