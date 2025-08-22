package biz

import (
	"github.com/google/wire"
	"review-service/internal/conf"
	"review-service/pkg/snowflake"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(NewReviewerUsecase, NewSnowflake)

// NewSnowflake 创建并返回雪花ID生成器
func NewSnowflake(c *conf.Data) (*snowflake.Snowflake, error) {
	return snowflake.NewSnowflake(c.Snowflake.WorkerID, c.Snowflake.DataCenterID)
}
