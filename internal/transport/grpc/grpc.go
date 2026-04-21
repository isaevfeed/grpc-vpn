package grpc

import (
	"fmt"
	"google.golang.org/grpc"
	"grpc-vpn/internal/transport/grpc/pb"
	"io"
	"log"
	"net"
	"sync"
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
	var targetConn net.Conn
	var once sync.Once
	for {
		packet, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return fmt.Errorf("stream.Recv: %w", err)
		}

		once.Do(func() {
			targetConn, err = net.Dial("tcp", packet.Target)
			if err != nil {
				log.Println("net.Dial", err)
				return
			}

			go func() {
				buf := make([]byte, 4096)

				for {
					n, err := targetConn.Read(buf)
					if err == io.EOF {
						log.Println("targetConn.Read", err)
						return
					}

					err = stream.Send(&pb.Package{
						Data: buf[:n],
					})
					if err != nil {
						log.Println("stream.Send", err)
						return
					}
				}
			}()
		})

		if targetConn == nil {
			continue
		}

		_, err = targetConn.Write(packet.Data)
		if err != nil {
			return fmt.Errorf("targetConn.Write: %w", err)
		}
	}
}
