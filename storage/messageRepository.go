package storage

import (
	"fmt"

	server "github.com/rossmerr/events-server/domain/server/v1"
	"go.etcd.io/etcd/server/v3/etcdserver"
	"go.etcd.io/etcd/server/v3/lease"
	"google.golang.org/protobuf/proto"
)

type MessageRepository struct {
	server *etcdserver.EtcdServer
}

func NewMessageRepository(server *etcdserver.EtcdServer) *MessageRepository {
	return &MessageRepository{
		server: server,
	}
}

func (s *MessageRepository) Save(msg *server.Message, stream *server.Stream, lease lease.LeaseID) (int64, error) {
	key := []byte(fmt.Sprintf("*ev.msg.%s.%s", stream.Id, msg.Id))
	bytes, err := proto.Marshal(msg)

	if err != nil {
		return 0, err
	}
	rev := s.server.KV().Put(key, bytes, lease)

	return rev, nil

}
