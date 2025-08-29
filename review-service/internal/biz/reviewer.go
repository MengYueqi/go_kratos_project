package biz

import (
	"context"
	v1 "review-service/api/review/v1"
	"review-service/internal/data/model"
	"review-service/pkg/snowflake"

	"github.com/go-kratos/kratos/v2/log"
)

var (
// ErrUserNotFound is user not found.
// ErrUserNotFound = errors.NotFound(v1.ErrorReason_USER_NOT_FOUND.String(), "user not found")
)

// Reviewer is a Reviewer model.
type Reviewer struct {
	Hello string
}

// ReviewerRepo is a Greater repo.
type ReviewerRepo interface {
	SaveReview(context.Context, *model.ReviewInfo) (*model.ReviewInfo, error)
	Update(context.Context, *Reviewer) (*Reviewer, error)
	GetReviewByOrderID(context.Context, int64) ([]*model.ReviewInfo, error)
	ListByHello(context.Context, string) ([]*Reviewer, error)
	ListAll(context.Context) ([]*Reviewer, error)
	DeleteReview(context.Context, int64) error
	GetReviewByID(context.Context, int64) (*model.ReviewInfo, error)
	GetReviewByReviewID(context.Context, int64) ([]*model.ReviewInfo, error)
	UpdateReviewByReviewID(context.Context, *model.ReviewInfo) (int64, error)
	GetReviewByUID(context.Context, int64) ([]*model.ReviewInfo, error)
	AddReviewReply(context.Context, *model.ReviewReplyInfo) (int64, error)
	AddAppealReview(context.Context, *model.ReviewAppealInfo) (int64, error)
	GetAppealByReviewID(context.Context, int64) ([]*model.ReviewAppealInfo, error)
}

// ReviewerUsecase is a Reviewer usecase.
type ReviewerUsecase struct {
	repo ReviewerRepo
	sf   *snowflake.Snowflake
	log  *log.Helper
}

// NewReviewerUsecase new a Reviewer usecase.
func NewReviewerUsecase(repo ReviewerRepo, sf *snowflake.Snowflake, logger log.Logger) *ReviewerUsecase {
	return &ReviewerUsecase{repo: repo, sf: sf, log: log.NewHelper(logger)}
}

// CreateReviewer creates a Reviewer, and returns the new Reviewer.
func (uc *ReviewerUsecase) CreateReviewer(ctx context.Context, review *model.ReviewInfo) (*model.ReviewInfo, error) {
	// 数据校验
	reviews, err := uc.repo.GetReviewByOrderID(ctx, review.OrderID)
	if err != nil {
		return nil, v1.ErrorDbFailed("DB search error!")
	}
	// 当前 Order 已经被评价
	if len(reviews) > 0 {
		return nil, v1.ErrorOrderReviewed("order id %d already exist a review", review.OrderID)
	}
	// 生成 ID
	// 使用雪花算法生成 ID
	review.ReviewID = uc.sf.NextID()
	uc.log.WithContext(ctx).Infof("[biz] CreateReviewer ID: %v", review.ReviewID)
	return uc.repo.SaveReview(ctx, review)
}

// 删除一个评论业务逻辑
func (uc *ReviewerUsecase) DeleteReviewer(ctx context.Context, ID int64) error {
	data, err := uc.repo.GetReviewByID(ctx, ID)
	if err != nil {
		return err
	}
	if data == nil {
		return v1.ErrorIdErr("Do not exist ID: %v", ID)
	} else if data.DeleteAt != nil {
		return v1.ErrorReviewHasBeenDeleted("Has been Delete ID: %v", ID)
	}

	return uc.repo.DeleteReview(ctx, ID)
}

// 根据 reviewID 获取评论内容
func (uc *ReviewerUsecase) GetReviewByReviewID(ctx context.Context, reviewId int64) (*model.ReviewInfo, error) {
	// 获取评论信息主逻辑
	info, err := uc.repo.GetReviewByReviewID(ctx, reviewId)
	if err != nil {
		return nil, err
	}
	return info[0], nil
}

// 根据 reviewId 更新数据
func (uc *ReviewerUsecase) UpdateReviewByReviewID(ctx context.Context, review *model.ReviewInfo) (int64, error) {
	rv, err := uc.repo.GetReviewByReviewID(ctx, review.ReviewID)
	if err != nil {
		return 0, err
	}
	if len(rv) == 0 {
		return 0, v1.ErrorReviewidErr("Do not exist ReviewID: %v", review.ReviewID)
	}
	if rv[0].DeleteAt != nil {
		return 0, v1.ErrorReviewidErr("The review has been delete: %v", review.ReviewID)
	}
	// 更新 review 主逻辑
	reviewId, err := uc.repo.UpdateReviewByReviewID(ctx, review)
	return reviewId, err
}

// 根据 uid 获取一个用户所有的评论
func (uc *ReviewerUsecase) ListReviewByUid(ctx context.Context, uid int64) ([]*model.ReviewInfo, error) {
	rvList, err := uc.repo.GetReviewByUID(ctx, uid)
	if err != nil {
		return nil, err
	}
	return rvList, nil
}

// 商家对用户的评论进行回复
func (uc *ReviewerUsecase) AddReplyReview(ctx context.Context, reply *model.ReviewReplyInfo) (int64, error) {
	// 生成雪花 ID
	reply.ReplyID = uc.sf.NextID()
	uc.log.WithContext(ctx).Infof("[biz] CreateReviewer ID: %v", reply.ReplyID)
	// 获取需要回复的评论
	reviewInfo, err := uc.repo.GetReviewByReviewID(ctx, reply.ReviewID)
	if err != nil {
		return 0, err
	}
	if len(reviewInfo) == 0 {
		return 0, v1.ErrorReviewidErr("Do not exist ReviewID: %v", reply.ReviewID)
	}
	if reviewInfo[0].ReviewID != reply.ReviewID {
		return 0, v1.ErrorStoreidReviewidMismatch("StoreID and Review's StoreID mismatch: %v - %v", reply.StoreID, reply.StoreID)
	}
	return uc.repo.AddReviewReply(ctx, reply)
}

// 商家对用户评论进行申诉
func (uc *ReviewerUsecase) AppealReview(ctx context.Context, appeal *model.ReviewAppealInfo) (int64, error) {
	// 1. 检查评论是否存在
	reviewInfo, err := uc.repo.GetReviewByReviewID(ctx, appeal.ReviewID)
	if err != nil {
		return 0, err
	}
	// 2. 检查评论是否已经被删除
	if len(reviewInfo) == 0 {
		return 0, v1.ErrorReviewidErr("Do not have ReviewID: %v", appeal.ReviewID)
	}
	// 3. 检查 StoreID 是否和评论的 StoreID 一致
	if reviewInfo[0].StoreID != appeal.StoreID {
		return 0, v1.ErrorStoreidReviewidMismatch("StoreID and Review's StoreID mismatch: %v - %v", appeal.StoreID, reviewInfo[0].StoreID)
	}
	// 4. 检查该评论是否已经被申诉过
	existAppeal, err := uc.repo.GetAppealByReviewID(ctx, appeal.ReviewID)
	if err != nil {
		return 0, err
	}
	if len(existAppeal) > 0 {
		return 0, v1.ErrorErrorAppealExists("The review has been appealed: %v", appeal.ReviewID)
	}

	// 主逻辑：插入一条申诉记录
	// 生成雪花 ID
	appeal.AppealID = uc.sf.NextID()
	uc.log.WithContext(ctx).Infof("[biz] CreateReviewer ID: %v", appeal.AppealID)
	review, err := uc.repo.AddAppealReview(ctx, appeal)
	if err != nil {
		return 0, err
	}
	return review, nil
}
