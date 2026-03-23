package grpc

import (
	"github.com/ThatSneakyCoder/RoutePulse/shared/env"
	pb "github.com/ThatSneakyCoder/RoutePulse/shared/proto/tracking"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type TrackingServiceClient struct {
	Client pb.TrackingServiceClient
	conn   *grpc.ClientConn
}

func NewTrackingServiceClient() (*TrackingServiceClient, error) {
	trackingServiceURL := env.GetString("TRACKING_SERVICE_ADDR", "tracking-service:9093")

	conn, err := grpc.NewClient(trackingServiceURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := pb.NewTrackingServiceClient(conn)

	return &TrackingServiceClient{
		Client: client,
		conn:   conn,
	}, nil
}

func (c *TrackingServiceClient) Close() {
	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			return
		}
	}
}
