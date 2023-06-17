package service

import (
	"github.com/hongliang5316/midjourney-apiserver/internal/common"
	"github.com/hongliang5316/midjourney-apiserver/pkg/api"
	"golang.org/x/sync/semaphore"
)

type Service struct {
	api.UnimplementedAPIServiceServer
	*common.Base
	*semaphore.Weighted
}

func New(base *common.Base) *Service {
	return &Service{
		Base: base,
	}
}
