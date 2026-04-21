package grpc

import (
	"fmt"
	"google.golang.org/grpc"
	"grpc-vpn/internal/transport/grpc/pb"
	"io"
)

type transport struct {
	pb.UnimplementedVPNServiceServer
}

func New() *grpc.Server {
	srv := grpc.NewServer()

	pb.RegisterVPNServiceServer(srv, &transport{})

	return srv
}

func (t *transport) Tunnel(stream pb.VPNService_TunnelServer) error {
	for {
		packet, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return fmt.Errorf("stream.Recv: %w", err)
		}
		fmt.Println("packet data len:", len(packet.Data))

		err = stream.Send(packet)
		if err != nil {
			return fmt.Errorf("stream.Send: %w", err)
		}
	}
}
