package data

import (
	"context"
	v1 "review-b/api/review/v1"

	"review-b/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type businessRepo struct {
	data *Data
	log  *log.Helper
}

// NewGreeterRepo .
func NewBusinessRepo(data *Data, logger log.Logger) biz.BusinessRepo {
	return &businessRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *businessRepo) Reply(ctx context.Context, replay *biz.ReplyParam) (int64, error) {
	r.log.WithContext(ctx).Infof("[data] Reply review Info: %+v", replay)
	// 调用 review 服务的 AddReplyReview 方法
	rId, err := r.data.rc.AddReplyReview(ctx, &v1.AddReplyReviewRequest{
		ReviewID:  replay.ReviewID,
		Content:   replay.Content,
		StoreID:   replay.StoreID,
		PicInfo:   replay.PicInfo,
		VideoInfo: replay.VideoInfo,
	})

	if err != nil {
		return -1, err
	}
	return rId.ReplyID, nil
}
