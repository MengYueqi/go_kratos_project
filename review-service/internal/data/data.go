package data

import (
	"context"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"review-service/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewGreeterRepo, NewDB, NewRedis)

// Data .
type Data struct {
	// TODO wrapped database client
	db    *gorm.DB
	redis *redis.Client
}

// NewData .
func NewData(db *gorm.DB, redis *redis.Client, c *conf.Data, logger log.Logger) (*Data, func(), error) {
	// 创建返回数据库连接实例
	dbInstance := &Data{
		db:    db,
		redis: redis,
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
