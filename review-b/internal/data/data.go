package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	v1 "review-b/api/review/v1"
	"review-b/internal/conf"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewBusinessRepo, NewReviewServiceClient)

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
func NewReviewServiceClient(c *conf.Data, logger log.Logger) v1.ReviewClient {
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint("127.0.0.1:9092"),
		grpc.WithTimeout(3600*time.Second),
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
