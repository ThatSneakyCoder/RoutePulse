package grpc

import (
	"os"

	pb "github.com/ThatSneakyCoder/RoutePulse/shared/proto/identity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type IdentityServiceClient struct {
	Client pb.IdentityServiceClient
	conn   *grpc.ClientConn
}

func NewIdentityServiceClient() (*IdentityServiceClient, error) {
	identityServiceURL := os.Getenv("IDENTITY_SERVICE_ADDR")
	if identityServiceURL == "" {
		identityServiceURL = "identity-service:9090"
	}

	conn, err := grpc.NewClient(identityServiceURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := pb.NewIdentityServiceClient(conn)

	return &IdentityServiceClient{
		Client: client,
		conn:   conn,
	}, nil
}

func (c *IdentityServiceClient) Close() {
	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			return
		}
	}
}
