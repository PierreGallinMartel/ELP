package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"image/jpeg"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	//imgPath := "C:/Users/pierr/Desktop/Programmation/ELP/ELP/GO/random_stuff/images/satellite.jpg"
	if len(os.Args) < 2 {
		fmt.Println("Please give the image path as an argument")
		os.Exit(1)
	}
	if len(os.Args) > 2 {
		fmt.Println("Too many arguments")
		os.Exit(1)
	}
	imgPath := os.Args[1]
	print("c'est" + string(imgPath) + "ahahah")
	c, err := net.Dial("tcp", ":8001")
	check(err)
	fmt.Println("Successfully connected to " + c.RemoteAddr().String())

	fmt.Fprintf(c, openImage(imgPath)+"\n") //When connected, send the image as a base64 string
	message, _ := bufio.NewReader(c).ReadString('\n')
	fmt.Print("->: " + message)
	byteImage, _ := base64.StdEncoding.DecodeString(message)
	jpgImage, err := jpeg.Decode(bytes.NewReader(byteImage))
	print("ahaha")
	check(err)
	ext := filepath.Ext(imgPath)
	name := strings.TrimSuffix(filepath.Base(imgPath), ext)
	newImagePath := fmt.Sprintf("%s/%s_blurred_edged_mthreaded%s", filepath.Dir(imgPath), name, ext)
	fg, err := os.Create(newImagePath)
	defer fg.Close()
	check(err)
	err = jpeg.Encode(fg, jpgImage, nil)
	print("done")
	check(err)
}

func openImage(imgPath string) string {

	f, err := os.Open(imgPath)
	check(err)
	defer f.Close()

	reader := bufio.NewReader(f)
	content, _ := ioutil.ReadAll(reader)

	// Encode as base64.
	encoded := base64.StdEncoding.EncodeToString(content)
	fmt.Println("ENCODED: " + strconv.Itoa(len(encoded)))
	return encoded
}
