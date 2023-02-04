Project Description:

This project implements multithreading in order to process heavy jpeg images for edge detection. This is performed through a client-server structure, which enables several clients to be served concurrently.
The processes are as follows:
- a Gaussian Blurr in order to smoothen the edges 
- a Laplacian Filter, which is highly sensitive to noise and other high frequency variations hence the necessity of pre-processing the image

In order to run the program, first run the server.go file. The server.go file takes arguments as follows:
	$go run server.go worker_quantity
e.g     $go run server.go 20
The worker_quantity is optional and defautls to 20 if no value is specified

Then the client.go file. The client.go file takes arguments as follows:
	$go run client.go image_path filtering_threshold
e.g     $go run client.go "C:/Users/pierr/Desktop/Programmation/ELP/GO/images/cathedral.jpg" 80
The processed image will be created in the folder of the source image. A threshold of 80 is recommanded for most images.

Some sample images have been provided for experimentation, the heavy one being satellite.jpg. You are welcome to feed heavier images to the program to better take advantage of the worker pools.
