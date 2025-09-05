package data

import (
	"context"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"review-service/internal/conf"
	"review-service/internal/data/query"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewReviewerRepo, NewDB, NewRedis, NewESClient)

// Data .
type Data struct {
	// TODO wrapped database client
	//db    *gorm.DB
	query *query.Query
	redis *redis.Client
	log   *log.Helper
	es    *elasticsearch.TypedClient
}

// NewData .
func NewData(db *gorm.DB, redis *redis.Client, esClient *elasticsearch.TypedClient, c *conf.Data, logger log.Logger) (*Data, func(), error) {
	// 为生成的代码制定对象
	query.SetDefault(db)

	// 创建返回数据库连接实例
	dbInstance := &Data{
		query: query.Q,
		redis: redis,
		log:   log.NewHelper(logger),
		es:    esClient,
	}

	// 关闭连接
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
		sqlDB, _ := db.DB()
		_ = sqlDB.Close()
		_ = redis.Close()
	}
	return dbInstance, cleanup, nil
}

// 连接 MySQL
func NewDB(c *conf.Data, logger log.Logger) (*gorm.DB, error) {
	helper := log.NewHelper(logger)
	db, err := gorm.Open(mysql.Open(c.Database.Source), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	helper.Infof("connect mysql success: %s", c.Database.Source)
	return db, nil
}

// 连接 Redis
// NewRedis creates a new redis.Client instance.
func NewRedis(c *conf.Data, logger log.Logger) (*redis.Client, error) {
	helper := log.NewHelper(logger)
	rdb := redis.NewClient(&redis.Options{
		Addr:     c.Redis.Addr,
		Password: c.Redis.Password,
		DB:       int(c.Redis.Db),
	})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}
	helper.Infof("connect redis success: %s", c.Redis.Addr)
	return rdb, nil
}

// ESClient 构造函数
func NewESClient(c *conf.Data, logger log.Logger) (*elasticsearch.TypedClient, error) {
	//helper := log.NewHelper(logger)
	// ES 配置
	cfg := elasticsearch.Config{
		Addresses: c.Elasticsearch.Addr,
	}

	return elasticsearch.NewTypedClient(cfg)
}
