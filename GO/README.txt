Project Description:

This project implements multithreading in order to process heavy jpeg images for edge detection. This is performed through a client-server structure, which enables several clients to be served concurrently.
The processes are as follows:
- a Gaussian Blurr in order to smoothen the edges 
- a Laplacian Filter, which is highly sensible to noise and other high frequency variations hence the necessity of pre-processing the image

In order to run the program, first run the server.go file, then the client.go file. The client.go file takes arguments as follows:
go run client.go image_path filtering_threshold (a value of 80 is recommended for most situations)