package services

import "github.com/google/uuid"

var (
	url   = []byte("github.com/rossmerr/events-server")
	Space = uuid.NewSHA1(uuid.NameSpaceURL, url)
)
