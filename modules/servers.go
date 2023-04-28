package modules

import (
	"github.com/samber/do"
	"google.golang.org/grpc"
)

// Wireup all services for the IOC container.
func RegisterServers(grpcServer *grpc.Server) error {

	do.New()

	return nil

}
