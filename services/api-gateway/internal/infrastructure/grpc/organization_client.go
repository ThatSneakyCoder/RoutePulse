package grpc

import (
	"github.com/ThatSneakyCoder/RoutePulse/shared/env"
	pb "github.com/ThatSneakyCoder/RoutePulse/shared/proto/organization"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type OrganizationServiceClient struct {
	Client pb.OrganizationServiceClient
	conn   *grpc.ClientConn
}

func NewOrganizationServiceClient() (*OrganizationServiceClient, error) {
	organizationServiceURL := env.GetString("ORGANIZATION_SERVICE_ADDR", "organization-service:9091")

	conn, err := grpc.NewClient(organizationServiceURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := pb.NewOrganizationServiceClient(conn)

	return &OrganizationServiceClient{
		Client: client,
		conn:   conn,
	}, nil
}

func (c *OrganizationServiceClient) Close() {
	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			return
		}
	}
}
