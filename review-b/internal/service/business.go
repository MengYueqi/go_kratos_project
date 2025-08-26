package service

import (
	"context"
	"review-b/internal/biz"

	pb "review-b/api/business/v1"
)

type BusinessService struct {
	pb.UnimplementedBusinessServer
	uc *biz.BusinessUsecase
}

func NewBusinessService(uc *biz.BusinessUsecase) *BusinessService {
	return &BusinessService{uc: uc}
}

func (s *BusinessService) ReplyUserReview(ctx context.Context, req *pb.ReplyReviewRequest) (*pb.ReplyReviewReply, error) {
	replyID, err := s.uc.CreateReply(ctx, &biz.ReplyParam{
		ReviewID:  req.GetReviewID(),
		Content:   req.GetContent(),
		StoreID:   req.GetStoreID(),
		PicInfo:   req.GetPicInfo(),
		VideoInfo: req.GetVideoInfo(),
	})
	if err != nil {
		return &pb.ReplyReviewReply{
			ReplyID: 0,
		}, err
	}
	return &pb.ReplyReviewReply{
		ReplyID: replyID,
	}, nil
}
