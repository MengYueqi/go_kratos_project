package service

import (
	"context"
	pb "review-service/api/review/v1"
	"review-service/internal/biz"
	"review-service/internal/data/model"
	"strconv"
	"time"
)

type ReviewService struct {
	pb.UnimplementedReviewServer
	uc *biz.ReviewerUsecase
}

func NewReviewService(uc *biz.ReviewerUsecase) *ReviewService {
	return &ReviewService{uc: uc}
}

// 创建一个评论
func (s *ReviewService) CreateReview(ctx context.Context, req *pb.CreateReviewRequest) (*pb.CreateReviewReply, error) {
	var anonymous int32
	if req.Anonymous {
		anonymous = 1
	}

	// 信息填入结构体
	data, err := s.uc.CreateReviewer(ctx, &model.ReviewInfo{
		CreateBy:     strconv.FormatInt(req.UserID, 10),
		UpdateBy:     strconv.FormatInt(req.UserID, 10),
		CreateAt:     time.Now(),
		UpdateAt:     time.Now(),
		Score:        req.Score,
		Status:       0,
		Anonymous:    anonymous,
		ServiceScore: req.ServiceScore,
		ExpressScore: req.ExpressScore,
		Content:      req.Content,
		PicInfo:      req.PicInfo,
		VideoInfo:    req.VideoInfo,
		OrderID:      req.OrderID,
	})
	// 错误处理
	if err != nil {
		panic(err)
	}
	// 返回数据
	return &pb.CreateReviewReply{
		ReviewID: data.ReviewID,
	}, nil
}
func (s *ReviewService) UpdateReview(ctx context.Context, req *pb.UpdateReviewRequest) (*pb.UpdateReviewReply, error) {
	return &pb.UpdateReviewReply{}, nil
}
func (s *ReviewService) DeleteReview(ctx context.Context, req *pb.DeleteReviewRequest) (*pb.DeleteReviewReply, error) {
	return &pb.DeleteReviewReply{}, nil
}
func (s *ReviewService) GetReview(ctx context.Context, req *pb.GetReviewRequest) (*pb.GetReviewReply, error) {
	return &pb.GetReviewReply{}, nil
}
func (s *ReviewService) ListReview(ctx context.Context, req *pb.ListReviewRequest) (*pb.ListReviewReply, error) {
	return &pb.ListReviewReply{}, nil
}
