package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"image/jpeg"
	"net"
	"os"
	"path/filepath"
	"strings"
)

const (
	con_host = "localhost"
	con_port = ":8000"
)

var connectedClients = 0
var received = ""

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func handleCo(c net.Conn) {
	fmt.Println("New client connected with " + c.RemoteAddr().String())

	for {
		temp, err := bufio.NewReader(c).ReadString('\n') // receive message
		check(err)

		if temp == "STOP" {
			fmt.Println("Client disconnected")
			break
		} else {
			byteImage, _ := base64.StdEncoding.DecodeString(temp)    //returns the bytes represented by the base64 string
			jpgImage, err := jpeg.Decode(bytes.NewReader(byteImage)) //Decodes the image using a byte reader
			check(err)

			// For now we just check if we receive the image correctly by exporting it. This is where your add your Gaussian blur magic
			ext := filepath.Ext("/Users/alicelebihan/Desktop/cathedral.jpg")
			name := strings.TrimSuffix(filepath.Base("/Users/alicelebihan/Desktop/cathedral.jpg"), ext)
			newImagePath := fmt.Sprintf("%s/%s_blurred%s", filepath.Dir("/Users/alicelebihan/Desktop/cathedral.jpg"), name, ext)
			fg, err := os.Create(newImagePath)
			defer fg.Close()
			check(err)
			err = jpeg.Encode(fg, jpgImage, nil)
			check(err)
			fmt.Println("Correctly receiveced image")
			//c.Write([]byte("Correctly received image\n"))
		}
	}
	c.Close()

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
