package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"strings"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	c, err := net.Dial("tcp", ":8000")
	check(err)
	fmt.Println("Successfully connected to " + c.RemoteAddr().String())

	fmt.Fprintf(c, openImage()+"\n") //When connected, send the image as a base64 string
	//the rest is just to have a way to communicate with the serer through the terminal and stop it
	for {
		//send what we write in the console
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(">> ")
		text, _ := reader.ReadString('\n')
		fmt.Fprintf(c, text+"\n")

		message, _ := bufio.NewReader(c).ReadString('\n')
		fmt.Print("->: " + message)
		if strings.TrimSpace(string(text)) == "STOP" {
			fmt.Println("TCP client exiting...")
			return
		}
	}
}

func openImage() string {
	imgPath := "/Users/alicelebihan/Desktop/cathedral.jpg"
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
