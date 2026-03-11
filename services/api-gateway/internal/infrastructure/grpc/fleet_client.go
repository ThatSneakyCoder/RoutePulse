package grpc

import (
	"github.com/ThatSneakyCoder/RoutePulse/shared/env"
	pb "github.com/ThatSneakyCoder/RoutePulse/shared/proto/fleet"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type FleetServiceClient struct {
	Client pb.FleetServiceClient
	conn   *grpc.ClientConn
}

func NewFleetServiceClient() (*FleetServiceClient, error) {
	fleetServiceURL := env.GetString("FLEET_SERVICE_ADDR", "fleet-service:9092")

	conn, err := grpc.NewClient(
		fleetServiceURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	client := pb.NewFleetServiceClient(conn)

	return &FleetServiceClient{
		Client: client,
		conn:   conn,
	}, nil
}

func (c *FleetServiceClient) Close() {
	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			return
		}
	}
}