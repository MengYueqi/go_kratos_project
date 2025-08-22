package biz

import (
	"context"
	"fmt"
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
		return nil, err
	}
	// 当前 Order 已经被评价
	if len(reviews) > 0 {
		return nil, fmt.Errorf("order id %d already exists", review.OrderID)
	}
	// 生成 ID
	// 使用雪花算法生成 ID
	review.ReviewID = uc.sf.NextID()
	uc.log.WithContext(ctx).Infof("[biz] CreateReviewer ID: %v", review.ID)
	return uc.repo.SaveReview(ctx, review)
}
