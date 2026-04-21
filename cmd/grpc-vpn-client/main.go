package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"grpc-vpn/internal/transport/grpc/pb"
	"log"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewVPNServiceClient(conn)

	stream, err := client.Tunnel(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	err = stream.Send(&pb.Package{
		Data:   []byte("GET / HTTP/1.1\r\nHost: google.com\r\n\r\n"),
		Target: "google.com:80",
	})
	if err != nil {
		log.Fatal(err)
	}

	resp, err := stream.Recv()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("tunnel response:", string(resp.Data))
}
