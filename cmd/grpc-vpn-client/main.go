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

	msg := []byte("Hello, world")

	err = stream.Send(&pb.Package{Data: msg})
	if err != nil {
		log.Fatal(err)
	}

	resp, err := stream.Recv()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("tunnel response:", string(resp.Data))
}
