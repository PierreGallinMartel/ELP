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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func convolution_Worker(wg *sync.WaitGroup, img *image.Image, wImg *image.RGBA, jobs <-chan [2]int, matrice [][]int) {
	var sumR float64
	var sumG float64
	var sumB float64
	var r float64
	var g float64
	var b float64
	var originalColor color.RGBA
	for true {
		cords, ok := <-jobs
		if ok == false {
			break
		}
		sumR = 0
		sumB = 0
		sumG = 0
		x := cords[0]
		y := cords[1]
		for i := -1; i <= 1; i++ {
			for j := -1; j <= 1; j++ {
				var imageX int
				var imageY int
				imageX = x + i
				imageY = y + j
				pixel := (*img).At(imageX, imageY)
				originalColor = color.RGBAModel.Convert(pixel).(color.RGBA)
				r = float64(originalColor.R)
				g = float64(originalColor.G)
				b = float64(originalColor.B)
				sumR = (sumR + (r * float64(matrice[i+1][j+1])))
				sumG = (sumG + (g * float64(matrice[i+1][j+1])))
				sumB = (sumB + (b * float64(matrice[i+1][j+1])))
			}
		}
		wImg.Set(cords[0], cords[1], color.NRGBA{
			uint8(sumR / 9),
			uint8(sumG / 9),
			uint8(sumB / 9),
			255,
		})
		wg.Done()
	}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	var wg sync.WaitGroup
	var matrice = [][]int{{1, 1, 1}, {1, 1, 1}, {1, 1, 1}}
	const numJobs = 1000
	jobs := make(chan [2]int, numJobs)
	imgPath := "C:/Users/pierr/Desktop/Programmation/ELP/ELP/GO/random_stuff/images/girl.jpg"
	f, err := os.Open(imgPath)
	check(err)
	defer f.Close()
	img, imType, imerr := image.Decode(f)
	print(imType)
	print(imerr)
	size := img.Bounds().Size()
	rect := image.Rect(0, 0, size.X, size.Y)
	wImg := image.NewRGBA(rect)
	for w := 0; w <= 8; w++ {
		go convolution_Worker(&wg, &img, wImg, jobs, matrice)
	}
	for x := 0; x < size.X; x++ {
		for y := 0; y < size.Y; y++ {
			wg.Add(1)
		}
	}
	go func() {
		for x := 0; x < size.X; x++ {
			for y := 0; y < size.Y; y++ {
				cords := [2]int{x, y}
				jobs <- cords

			}
		}
	}()
	wg.Wait()
	close(jobs)
	ext := filepath.Ext(imgPath)
	name := strings.TrimSuffix(filepath.Base(imgPath), ext)
	newImagePath := fmt.Sprintf("%s/%s_goblurred%s", filepath.Dir(imgPath), name, ext)
	fg, err := os.Create(newImagePath)
	defer fg.Close()
	check(err)
	err = jpeg.Encode(fg, wImg, nil)
	check(err)

}
