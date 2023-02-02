package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"math"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	con_host = "localhost"
	con_port = ":8000"
)

var connectedClients = 0
var received = ""

func TrimSuffix(s, suffix string) string {
	if strings.HasSuffix(s, suffix) {
		s = s[:len(s)-len(suffix)]
	}
	return s
}
func handleCo(c net.Conn) {
	fmt.Println("New client connected with " + c.RemoteAddr().String())
	temp, err := bufio.NewReader(c).ReadString('\n') // receive message
	check(err)
	s_string, err := bufio.NewReader(c).ReadString('\n')
	s_string = TrimSuffix(s_string, "\n")
	s, _ := strconv.Atoi(s_string)
	check(err)
	byteImage, _ := base64.StdEncoding.DecodeString(temp)    //returns the bytes represented by the base64 string
	jpgImage, err := jpeg.Decode(bytes.NewReader(byteImage)) //Decodes the image using a byte reader
	check(err)
	start := time.Now()
	result := imageProcessing(jpgImage, int(s))
	elapsed := time.Since(start)
	// For now we just check if we receive the image correctly by exporting it. This is where your add your Gaussian blur magic
	check(err)
	fmt.Println("processing took ", elapsed)
	fmt.Println("Correctly receiveced image")

	fmt.Fprintf(c, openImage(result)+"\n") //When connected, send the image as a base64 string
	//c.Write([]byte("Correctly received image\n"))

	c.Close()
	fmt.Println("Closed connection")
}

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
	for cords := range jobs {
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

func edgeDetection_Worker(wg *sync.WaitGroup, img *image.RGBA, wImg *image.RGBA, jobs <-chan [2]int, threshold float64, size image.Point) {
	var N float64
	var S float64
	var E float64
	var W float64
	var laplacien float64
	var I_final uint8
	for cords := range jobs {
		x := cords[0]
		y := cords[1]
		pixel := img.At(x, y)
		originalColor := color.RGBAModel.Convert(pixel).(color.RGBA)
		C := float64((originalColor.R + originalColor.G + originalColor.B) / (3))
		if y > 1 {
			pixel2 := img.At(x, y-1)
			originalColor2 := color.RGBAModel.Convert(pixel2).(color.RGBA)
			N = float64((originalColor2.R + originalColor2.G + originalColor2.B) / 3)
		} else {
			N = C
		}
		if y < size.Y-1 {
			pixel2 := img.At(x, y+1)
			originalColor2 := color.RGBAModel.Convert(pixel2).(color.RGBA)
			S = float64((originalColor2.R + originalColor2.G + originalColor2.B) / 3)
		} else {
			S = C
		}
		CV := N + S - 2*C
		if x > 1 {
			pixel2 := img.At(x-1, y)
			originalColor2 := color.RGBAModel.Convert(pixel2).(color.RGBA)
			W = float64((originalColor2.R + originalColor2.G + originalColor2.B) / 3)
		} else {
			W = C
		}
		if x < size.X-1 {
			pixel2 := img.At(x+1, y)
			originalColor2 := color.RGBAModel.Convert(pixel2).(color.RGBA)
			E = float64((originalColor2.R + originalColor2.G + originalColor2.B) / 3)
		} else {
			E = C
		}
		CH := W + E - 2*C
		laplacien = (CV + CH)
		if math.Abs(laplacien) < threshold {
			I_final = 0
		} else {
			I_final = 255
		}
		c := color.RGBA{
			R: I_final, G: I_final, B: I_final, A: 1,
		}
		wImg.Set(x, y, c)
		wg.Done()
	}

}

func imageProcessing(img image.Image, s_int int) *image.RGBA {
	worker_amount := 20
	if len(os.Args) > 1 {
		worker_amount_string := os.Args[1]
		worker_amount, _ = strconv.Atoi(worker_amount_string)
	}
	var wg sync.WaitGroup
	var matrice = [][]int{{1, 1, 1}, {1, 1, 1}, {1, 1, 1}}
	const numJobs = 1000
	jobs := make(chan [2]int, numJobs)
	size := img.Bounds().Size()
	rect := image.Rect(0, 0, size.X, size.Y)
	convolutionImage := image.NewRGBA(rect)
	edgeDetectionImage := image.NewRGBA(rect)
	fmt.Printf("with " + fmt.Sprint(worker_amount) + " workers: ")
	for w := 0; w < worker_amount; w++ {
		go convolution_Worker(&wg, &img, convolutionImage, jobs, matrice)
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
	jobs = make(chan [2]int, numJobs)
	var s float64
	s = float64(s_int)
	for w := 0; w < worker_amount; w++ {
		go edgeDetection_Worker(&wg, convolutionImage, edgeDetectionImage, jobs, s, size)
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
	return (edgeDetectionImage)
}

func main() {

	l, err := net.Listen("tcp", con_port)
	check(err)
	fmt.Println("Successfully started server on port " + con_port)

	defer l.Close()

	for {
		//accept client connection
		c, err := l.Accept()
		check(err)
		go handleCo(c) //allows us to accept several clients
		connectedClients++
	}
}

func openImage(result *image.RGBA) string {
	buf := new(bytes.Buffer)
	jpeg.Encode(buf, result, nil)
	// Encode as base64.
	blabla := buf.String()
	encoded := base64.StdEncoding.EncodeToString([]byte(blabla))
	//fmt.Printf(encoded)
	fmt.Println("ENCODED: " + strconv.Itoa(len(encoded)))
	return encoded
}
