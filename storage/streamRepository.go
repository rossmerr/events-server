package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	server "github.com/rossmerr/events-server/domain/server/v1"
	log "github.com/sirupsen/logrus"
	pb "go.etcd.io/etcd/api/v3/etcdserverpb"
	"go.etcd.io/etcd/server/v3/etcdserver"
	"go.etcd.io/etcd/server/v3/lease"
	"go.etcd.io/etcd/server/v3/mvcc"
	"google.golang.org/protobuf/proto"
)

type StreamRepository struct {
	server *etcdserver.EtcdServer
}

func NewStreamRepository(server *etcdserver.EtcdServer) *StreamRepository {
	return &StreamRepository{
		server: server,
	}
}

func (s *StreamRepository) Find(subject string) []*server.Stream {
	streams := []*server.Stream{}

	key := []byte(fmt.Sprintf("*ev.streams.%s", subject))
	results, err := s.server.KV().Range(context.Background(), key, key, mvcc.RangeOptions{})

	if err != nil {
		log.Error(errors.Join(fmt.Errorf("StreamRepository.Find"), err))
	}

	for _, kv := range results.KVs {
		stream := &server.Stream{}
		err := proto.Unmarshal(kv.Value, stream)
		if err != nil {
			log.Error(errors.Join(fmt.Errorf("StreamRepository.Find"), err))
		}

		streams = append(streams, stream)
	}

	return streams
}

func (s *StreamRepository) Save(stream *server.Stream) (string, error) {
	l, err := s.server.LeaseGrant(context.Background(), &pb.LeaseGrantRequest{
		TTL: 0,
	})

	if err != nil {
		return "", errors.Join(fmt.Errorf("domainServer.Publish"), err)
	}

	stream.Id = uuid.NewSHA1(Space, []byte(stream.Name)).String()
	key := []byte(fmt.Sprintf("*ev.streams.%s.%s", stream.Subject, stream.Id))

	bytes, err := proto.Marshal(stream)
	if err != nil {
		return "", errors.Join(fmt.Errorf("domainServer.Publish"), err)
	}

	s.server.KV().Put(key, bytes, lease.LeaseID(l.ID))
	return stream.Id, err
}

func (s *StreamRepository) Watch(subject string) (<-chan *server.Message, chan<- bool) {
	watchStream := s.server.KV().NewWatchStream()

	key := []byte(subject)

	_, err := watchStream.Watch(mvcc.WatchID(0), key, key, 0)
	if err != nil {
		log.Error(errors.Join(fmt.Errorf("domainServer.Subscribe"), err))
	}

	response := make(chan *server.Message)
	signals := make(chan bool)

	go func() {
		defer watchStream.Close()

		select {
		case watch := <-watchStream.Chan():
			for _, event := range watch.Events {

				var msg *server.Message
				if err := proto.Unmarshal(event.Kv.Value, msg); err != nil {
					log.Error(errors.Join(fmt.Errorf("domainServer.Subscribe"), err))
				} else {
					response <- msg
				}
			}
		case <-signals:
			break
		}

		close(signals)
		close(response)
	}()

	return response, signals
}
