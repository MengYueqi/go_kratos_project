package service

import (
	"context"
	"google.golang.org/protobuf/types/known/timestamppb"
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
		UserID:       req.UserID,
		StoreID:      req.StoreID,
	})
	// 错误处理
	if err != nil {
		return &pb.CreateReviewReply{}, err
	}
	// 返回数据
	return &pb.CreateReviewReply{
		ReviewID: data.ReviewID,
	}, nil
}

func (s *ReviewService) UpdateReview(ctx context.Context, req *pb.UpdateReviewRequest) (*pb.UpdateReviewReply, error) {
	var anonymous int32
	if req.Anonymous {
		anonymous = 1
	}
	newRV := &model.ReviewInfo{
		ReviewID:     req.ReviewID,
		Content:      req.Content,
		Score:        req.Score,
		ServiceScore: req.ServiceScore,
		ExpressScore: req.ExpressScore,
		PicInfo:      req.PicInfo,
		VideoInfo:    req.VideoInfo,
		Anonymous:    anonymous,
	}
	reviewId, err := s.uc.UpdateReviewByReviewID(ctx, newRV)
	if err != nil {
		return &pb.UpdateReviewReply{}, err
	}
	return &pb.UpdateReviewReply{
		ReviewID: reviewId,
	}, nil
}
func (s *ReviewService) DeleteReview(ctx context.Context, req *pb.DeleteReviewRequest) (*pb.DeleteReviewReply, error) {
	err := s.uc.DeleteReviewer(ctx, req.ID)
	if err != nil {
		return &pb.DeleteReviewReply{}, err
	}
	return &pb.DeleteReviewReply{
		ReviewID: req.ReviewID,
	}, nil
}
func (s *ReviewService) GetReview(ctx context.Context, req *pb.GetReviewRequest) (*pb.GetReviewReply, error) {
	rv, err := s.uc.GetReviewByReviewID(ctx, req.ReviewID)
	if err != nil {
		return &pb.GetReviewReply{}, err
	}
	var anonymous bool
	if rv.Anonymous == 1 {
		anonymous = true
	}
	return &pb.GetReviewReply{
		UserID:       rv.UserID,
		OrderID:      rv.OrderID,
		Score:        rv.Score,
		ServiceScore: rv.ServiceScore,
		ExpressScore: rv.ExpressScore,
		Content:      rv.Content,
		PicInfo:      rv.PicInfo,
		VideoInfo:    rv.VideoInfo,
		Anonymous:    anonymous,
		CreateTime:   timestamppb.New(rv.CreateAt),
		UpdateTime:   timestamppb.New(rv.UpdateAt),
	}, nil
}
func (s *ReviewService) ListReviewByUid(ctx context.Context, req *pb.ListReviewByUidRequest) (*pb.ListReviewByUidReply, error) {
	rvList, err := s.uc.ListReviewByUid(ctx, req.UserID)
	if err != nil {
		return &pb.ListReviewByUidReply{}, err
	}
	var retReviewList []*pb.ReviewReply
	for _, rv := range rvList {
		var anonymous bool
		if rv.Anonymous == 1 {
			anonymous = true
		}
		retReviewList = append(retReviewList, &pb.ReviewReply{
			UserID:       rv.UserID,
			OrderID:      rv.OrderID,
			Score:        rv.Score,
			ServiceScore: rv.ServiceScore,
			ExpressScore: rv.ExpressScore,
			Content:      rv.Content,
			PicInfo:      rv.PicInfo,
			VideoInfo:    rv.VideoInfo,
			Anonymous:    anonymous,
			CreateTime:   timestamppb.New(rv.CreateAt),
			UpdateTime:   timestamppb.New(rv.UpdateAt),
		})
	}
	return &pb.ListReviewByUidReply{
		Reviews: retReviewList,
	}, nil
}

func (s *ReviewService) AddReplyReview(ctx context.Context, req *pb.AddReplyReviewRequest) (*pb.AddReplyReviewReply, error) {
	replyInfo := &model.ReviewReplyInfo{
		CreateBy:  strconv.FormatInt(req.StoreID, 10),
		UpdateBy:  strconv.FormatInt(req.StoreID, 10),
		CreateAt:  time.Now(),
		UpdateAt:  time.Now(),
		StoreID:   req.StoreID,
		ReviewID:  req.ReviewID,
		Content:   req.Content,
		VideoInfo: req.VideoInfo,
		PicInfo:   req.PicInfo,
	}
	replyID, err := s.uc.AddReplyReview(ctx, replyInfo)
	if err != nil {
		return &pb.AddReplyReviewReply{}, err
	}
	return &pb.AddReplyReviewReply{
		ReplyID: replyID,
	}, nil
}

func (s *ReviewService) AppealReview(ctx context.Context, req *pb.AppealReviewRequest) (*pb.AppealReviewReply, error) {
	appealID, err := s.uc.AppealReview(ctx, &model.ReviewAppealInfo{
		ReviewID:  req.ReviewID,
		StoreID:   req.StoreID,
		Content:   req.Content,
		PicInfo:   req.PicInfo,
		VideoInfo: req.VideoInfo,
		Status:    10,
		CreateAt:  time.Now(),
		UpdateAt:  time.Now(),
		CreateBy:  strconv.FormatInt(req.StoreID, 10),
		UpdateBy:  strconv.FormatInt(req.StoreID, 10),
	})
	if err != nil {
		return &pb.AppealReviewReply{}, err
	}
	return &pb.AppealReviewReply{
		AppealID: appealID,
	}, nil
}

// O 端处理申述
func (s *ReviewService) HandleAppeal(ctx context.Context, req *pb.AppealOperateRequest) (*pb.AppealOperateReply, error) {
	data, err := s.uc.HandleAppeal(ctx, &model.ReviewAppealInfo{
		AppealID: req.AppealID,
		Status:   req.Status,
		ID:       req.ID,
		Reason:   req.Reason,
		OpUser:   req.OpUser,
		UpdateAt: time.Now(),
	})
	if err != nil {
		return &pb.AppealOperateReply{}, err
	}
	return &pb.AppealOperateReply{
		AppealID: data.AppealID,
		Status:   data.Status,
		ID:       data.ID,
		Reason:   data.Reason,
		OpUser:   data.OpUser,
	}, nil
}

// 根据 StoreID 查找评论
func (s *ReviewService) ListReviewByStoreID(ctx context.Context, req *pb.ListReviewByStoreIDRequest) (*pb.ListReviewByStoreIDReply, error) {
	data, err := s.uc.ListReviewByStoreID(ctx, req.StoreID, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}
	// 格式化
	list := make([]*pb.ReviewInfo, 0, len(data))
	for _, item := range data {
		var anonymous bool
		if item.Anonymous == 1 {
			anonymous = true
		}
		list = append(list, &pb.ReviewInfo{
			UserID:       item.UserID,
			OrderID:      item.OrderID,
			Score:        item.Score,
			ServiceScore: item.ServiceScore,
			ExpressScore: item.ExpressScore,
			Content:      item.Content,
			PicInfo:      item.PicInfo,
			VideoInfo:    item.VideoInfo,
			Anonymous:    anonymous,
			CreateTime:   timestamppb.New(time.Time(item.CreateAt)),
			UpdateTime:   timestamppb.New(time.Time(item.UpdateAt)),
		})
	}
	return &pb.ListReviewByStoreIDReply{
		Reviews: list,
	}, nil
}
