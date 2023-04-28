package services

import (
	"context"
	"errors"
	"fmt"
	"io"

	domain "github.com/rossmerr/events-server/domain/events/v1"
	log "github.com/sirupsen/logrus"
	pb "go.etcd.io/etcd/api/v3/etcdserverpb"
	"go.etcd.io/etcd/server/v3/etcdserver"
	"go.etcd.io/etcd/server/v3/lease"
	"go.etcd.io/etcd/server/v3/mvcc"
)

type domainServer struct {
	domain.UnimplementedDomainEventsServiceServer
	server *etcdserver.EtcdServer
}

func NewDomainServer(server *etcdserver.EtcdServer) domain.DomainEventsServiceServer {
	return &domainServer{
		server: server,
	}
}

func (s *domainServer) Publish(stream domain.DomainEventsService_PublishServer) error {
	log.Debug(fmt.Sprintf("domainServer.Publish  %+v", stream))

	for {
		publish, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&domain.PublishResponse{})
		}

		if err != nil {
			return errors.Join(fmt.Errorf("domainServer.Publish"), err)
		}

		l, err := s.server.LeaseGrant(context.Background(), &pb.LeaseGrantRequest{
			TTL: 0,
		})

		if err != nil {
			return errors.Join(fmt.Errorf("domainServer.Publish"), err)
		}

		key := []byte(publish.Subject)
		s.server.KV().Put(key, []byte(""), lease.LeaseID(l.ID))
	}

}

func (s *domainServer) Subscribe(request *domain.SubscribeRequest, stream domain.DomainEventsService_SubscribeServer) error {
	log.Debug(fmt.Sprintf("domainServer.Subscribe  %+v", request))

	watchStream := s.server.KV().NewWatchStream()
	defer watchStream.Close()

	key := []byte(request.Subject)

	_, err := watchStream.Watch(mvcc.WatchID(0), key, key, 0)
	if err != nil {
		return errors.Join(fmt.Errorf("domainServer.Subscribe"), err)
	}

	response := <-watchStream.Chan()

	for _, event := range response.Events {
		subject := string(event.Kv.Key)

		message := &domain.SubscribeResponse{
			Subject: subject,
		}

		if err := stream.Send(message); err != nil {
			return errors.Join(fmt.Errorf("domainServer.Subscribe"), err)
		}
	}

	return nil
}
