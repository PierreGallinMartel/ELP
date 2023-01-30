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
	"time"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Please give both the image path and the filtering threshold (recommended value: 80")
		os.Exit(1)
	}
	if len(os.Args) > 3 {
		fmt.Println("Too many arguments")
		os.Exit(1)
	}
	imgPath := os.Args[1]
	s_string := os.Args[2]
	c, err := net.Dial("tcp", ":8004")
	check(err)
	fmt.Println("Successfully connected to " + c.RemoteAddr().String())

	fmt.Fprintf(c, openImage(imgPath)+"\n")
	time.Sleep(1 * time.Second)
	fmt.Fprintf(c, s_string+"\n")
	message, _ := bufio.NewReader(c).ReadString('\n')
	fmt.Println("->: image came back")
	byteImage, _ := base64.StdEncoding.DecodeString(message)
	jpgImage, err := jpeg.Decode(bytes.NewReader(byteImage))
	check(err)
	ext := filepath.Ext(imgPath)
	name := strings.TrimSuffix(filepath.Base(imgPath), ext)
	newImagePath := fmt.Sprintf("%s/%s_filtered_at_%s %s", filepath.Dir(imgPath), name, s_string, ext)
	fg, err := os.Create(newImagePath)
	defer fg.Close()
	check(err)
	err = jpeg.Encode(fg, jpgImage, nil)
	fmt.Println("done")
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
