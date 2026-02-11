package grpc

import (
	"github.com/ThatSneakyCoder/RoutePulse/shared/env"
	pb "github.com/ThatSneakyCoder/RoutePulse/shared/proto/analytics"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AnalyticsServiceClient struct {
	Client pb.AnalyticsServiceClient
	conn   *grpc.ClientConn
}

func NewAnalyticsServiceClient() (*AnalyticsServiceClient, error) {
	analyticsServiceURL := env.GetString("ANALYTICS_SERVICE_ADDR", "analytics-service:9096")

	conn, err := grpc.NewClient(analyticsServiceURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := pb.NewAnalyticsServiceClient(conn)

	return &AnalyticsServiceClient{
		Client: client,
		conn:   conn,
	}, nil
}

func (c *AnalyticsServiceClient) Close() {
	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			return
		}
	}
}
