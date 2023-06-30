package services

import (
	"context"
	"errors"
	"fmt"
	"io"

	domain "github.com/rossmerr/events-server/domain/events/v1"
	server "github.com/rossmerr/events-server/domain/server/v1"
	"github.com/rossmerr/events-server/storage"

	log "github.com/sirupsen/logrus"
	pb "go.etcd.io/etcd/api/v3/etcdserverpb"
	"go.etcd.io/etcd/server/v3/etcdserver"
	"go.etcd.io/etcd/server/v3/lease"
)

type domainServer struct {
	domain.UnimplementedDomainEventsServiceServer
	server      *etcdserver.EtcdServer
	messageRepo storage.MessageRepository
	streamRepo  storage.StreamRepository
}

func NewDomainServer(server *etcdserver.EtcdServer) domain.DomainEventsServiceServer {
	return &domainServer{
		server: server,
	}
}
func (s *domainServer) Stream(ctx context.Context, request *domain.StreamRequest) (*domain.StreamResponse, error) {

	log.Debug(fmt.Sprintf("domainServer.Stream  %+v", request.Subject))

	response := &domain.StreamResponse{}

	id, err := s.streamRepo.Save(&server.Stream{
		Name:       request.Name,
		Subject:    request.Subject,
		Replicas:   request.Replicas,
		MaxDeliver: request.MaxDeliver,
	})

	response.Id = id

	return response, err

}

func (s *domainServer) Publish(stream domain.DomainEventsService_PublishServer) error {
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

		lease := lease.LeaseID(l.ID)

		if err != nil {
			return errors.Join(fmt.Errorf("domainServer.Publish"), err)
		}

		log.Debug(fmt.Sprintf("domainServer.Publish  %+v", publish.Subject))

		streams := s.streamRepo.Find(publish.Subject)

		msg := &server.Message{
			Id:      publish.Id,
			Subject: publish.Subject,
			Message: publish.Message,
		}

		for _, stream := range streams {
			_, err := s.messageRepo.Save(msg, stream, lease)
			if err != nil {
				return errors.Join(fmt.Errorf("domainServer.Publish"), err)
			}
		}

	}

}

func (s *domainServer) Subscribe(request *domain.SubscribeRequest, stream domain.DomainEventsService_SubscribeServer) error {
	log.Debug(fmt.Sprintf("domainServer.Subscribe  %+v", request.Subject))

	watchStream, signals := s.streamRepo.Watch(request.Subject)

	for msg := range watchStream {

		message := &domain.SubscribeResponse{
			Id:      msg.Id,
			Message: msg.Message,
		}

		if err := stream.Send(message); err != nil {

			if err != nil {
				signals <- true
			}

			if err == io.EOF {
				break
			}

			return errors.Join(fmt.Errorf("domainServer.Subscribe"), err)
		}
	}

	return nil

}
