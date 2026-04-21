package grpc

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"grpc-vpn/internal/transport/grpc/pb"
	"io"
	"log"
	"net"
)

type Client struct {
	conn   *grpc.ClientConn
	client pb.VPNServiceClient
}

func NewClient(addr string) (*Client, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("grpc.Dial: %w", err)
	}

	return &Client{
		conn:   conn,
		client: pb.NewVPNServiceClient(conn),
	}, nil
}

func (c *Client) Handle(ctx context.Context, conn net.Conn, target string) error {
	defer conn.Close()

	log.Println("TARGET:", target)

	stream, err := c.client.Tunnel(ctx)
	if err != nil {
		return fmt.Errorf("client.Tunnel: %w", err)
	}

	err = stream.Send(&pb.Package{
		Target: target,
	})
	if err != nil {
		return fmt.Errorf("stream.Send: %w", err)
	}

	var eg errgroup.Group
	eg.Go(func() error {
		buf := make([]byte, 4096)

		for {
			n, err := conn.Read(buf)
			if err != nil {
				stream.CloseSend()
				return fmt.Errorf("conn.Read: %w", err)
			}

			err = stream.Send(&pb.Package{
				Data: buf[:n],
			})
			if err != nil {
				stream.CloseSend()
				return fmt.Errorf("stream.Send: %w", err)
			}
		}
	})
	eg.Go(func() error {
		for {
			packet, err := stream.Recv()
			if err == io.EOF {
				return nil
			}
			if err != nil {
				return fmt.Errorf("stream.Recv: %w", err)
			}

			_, err = conn.Write(packet.Data)
			if err != nil {
				return fmt.Errorf("conn.Write: %w", err)
			}
		}
	})
	if err := eg.Wait(); err != nil {
		return fmt.Errorf("eg.Wait: %w", err)
	}

	return nil
}
