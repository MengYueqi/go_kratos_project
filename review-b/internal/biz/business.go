package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

var (
// ErrUserNotFound is user not found.
// ErrUserNotFound = errors.NotFound(v1.ErrorReason_USER_NOT_FOUND.String(), "user not found")
)

// Greeter is a Greeter model.
type Business struct {
	Hello string
}

type ReplyParam struct {
	ReviewID  int64
	Content   string
	StoreID   int64
	PicInfo   string
	VideoInfo string
}

type AppealParam struct {
	ReviewID  int64
	StoreID   int64
	Content   string
	PicInfo   string
	VideoInfo string
}

// BusinessRepo is a Business repo.
type BusinessRepo interface {
	Reply(ctx context.Context, param *ReplyParam) (int64, error)
	Appeal(ctx context.Context, param *AppealParam) (int64, error)
}

// GreeterUsecase is a Greeter usecase.
type BusinessUsecase struct {
	repo BusinessRepo
	log  *log.Helper
}

// NewGreeterUsecase new a Greeter usecase.
func NewBusinessUsecase(repo BusinessRepo, logger log.Logger) *BusinessUsecase {
	return &BusinessUsecase{repo: repo, log: log.NewHelper(logger)}
}

// CreateGreeter creates a Greeter, and returns the new Greeter.
func (uc *BusinessUsecase) CreateReply(ctx context.Context, r *ReplyParam) (int64, error) {
	uc.log.WithContext(ctx).Infof("CreateReply")
	return uc.repo.Reply(ctx, r)
}

// 添加申诉
func (uc *BusinessUsecase) AppealUserReview(ctx context.Context, appeal *AppealParam) (int64, error) {
	idx, err := uc.repo.Appeal(ctx, appeal)
	if err != nil {
		return 0, err
	}
	uc.log.WithContext(ctx).Infof("AppealUserReview")
	return idx, nil
}
