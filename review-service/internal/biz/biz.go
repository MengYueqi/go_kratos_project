package biz

import (
	"github.com/google/wire"
	"review-service/internal/conf"
	"review-service/internal/data/model"
	"review-service/pkg/snowflake"
	"strings"
	"time"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(NewReviewerUsecase, NewSnowflake)

// NewSnowflake 创建并返回雪花ID生成器
func NewSnowflake(c *conf.Data) (*snowflake.Snowflake, error) {
	return snowflake.NewSnowflake(c.Snowflake.WorkerID, c.Snowflake.DataCenterID)
}

type MyReviewInfo struct {
	*model.ReviewInfo
	CreateAt     MyTime `json:"create_at"`
	UpdateAt     MyTime `json:"update_at"`
	ID           int64  `json:"id,string"`            // 主键
	Version      int32  `json:"version,string"`       // 乐观锁标记
	ReviewID     int64  `json:"review_id,string"`     // 评价id
	Score        int32  `json:"score,string"`         // 评分
	ServiceScore int32  `json:"service_score,string"` // 商家服务评分
	ExpressScore int32  `json:"express_score,string"` // 物流评分
	HasMedia     int32  `json:"has_media,string"`     // 是否有图或视频
	OrderID      int64  `json:"order_id,string"`      // 订单id
	SkuID        int64  `json:"sku_id,string"`        // sku id
	SpuID        int64  `json:"spu_id,string"`        // spu id
	StoreID      int64  `json:"store_id,string"`      // 店铺id
	UserID       int64  `json:"user_id,string"`       // ⽤户id
	Anonymous    int32  `json:"anonymous,string"`     // 是否匿名
	Status       int32  `json:"status,string"`
	IsDefault    int32  `json:"is_default,string"` // 是否默认评价
	HasReply     int32  `json:"has_reply,string"`  // 是否有商家回复:0⽆;1有
}

type MyTime time.Time

// MarshalJSON 重写时间格式，并将其保存在 MyTime 结构体中
func (t *MyTime) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), `"`)
	tt, err := time.Parse(time.DateTime, s)
	if err != nil {
		return err
	}
	*t = MyTime(tt)
	return nil
}
