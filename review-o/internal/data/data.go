package data

import (
	"context"
	"github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/google/wire"
	"github.com/hashicorp/consul/api"
	v1 "review-o/api/review/v1"
	"review-o/internal/conf"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewOperationRepo, NewReviewServiceClient, NewDiscovery)

// Data .
type Data struct {
	// TODO wrapped database client
	rc v1.ReviewClient
}

// NewData
func NewData(c *conf.Data, rc v1.ReviewClient, logger log.Logger) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	return &Data{
		rc: rc,
	}, cleanup, nil
}

// 生成连接中台的客户端，用于 b 端 review 服务
func NewReviewServiceClient(discovery registry.Discovery, logger log.Logger) v1.ReviewClient {
	endpoint := "discovery:///review-service"
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint(endpoint),
		grpc.WithDiscovery(discovery),
		grpc.WithMiddleware(
			recovery.Recovery(),
			validate.Validator(),
		),
	)
	if err != nil {
		panic(err)
	}
	return v1.NewReviewClient(conn)
}

// 服务发现构造函数
func NewDiscovery(conf *conf.Registry) registry.Discovery {
	// new consul client
	c := api.DefaultConfig()
	c.Address = conf.Consul.Addr
	c.Scheme = conf.Consul.Scheme
	client, err := api.NewClient(c)
	if err != nil {
		panic(err)
	}
	// new dis with consul client
	dis := consul.New(client)

	return dis
}
