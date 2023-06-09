package service

import (
	"github.com/hongliang5316/midjourney-apiserver/internal/common"
	"github.com/hongliang5316/midjourney-apiserver/pkg/api"
)

type Service struct {
	api.UnimplementedAPIServiceServer
	*common.Base
}

func New(base *common.Base) *Service {
	return &Service{
		Base: base,
	}
}
