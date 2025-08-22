package data

import (
	"context"
	"review-service/internal/data/model"

	"review-service/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type ReviewerRepo struct {
	data *Data
	log  *log.Helper
}

// NewGreeterRepo .
func NewReviewerRepo(data *Data, logger log.Logger) biz.ReviewerRepo {
	return &ReviewerRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *ReviewerRepo) SaveReview(ctx context.Context, review *model.ReviewInfo) (*model.ReviewInfo, error) {
	err := r.data.query.ReviewInfo.WithContext(ctx).Save(review)
	return review, err
}

func (r *ReviewerRepo) Update(ctx context.Context, g *biz.Reviewer) (*biz.Reviewer, error) {
	return g, nil
}

func (r *ReviewerRepo) GetReviewByOrderID(ctx context.Context, orderId int64) ([]*model.ReviewInfo, error) {
	find, err := r.data.query.ReviewInfo.WithContext(ctx).Where(r.data.query.ReviewInfo.ReviewID.Eq(orderId)).Find()
	if err != nil {
		return nil, err
	}
	return find, nil
}

func (r *ReviewerRepo) ListByHello(context.Context, string) ([]*biz.Reviewer, error) {
	return nil, nil
}

func (r *ReviewerRepo) ListAll(context.Context) ([]*biz.Reviewer, error) {
	return nil, nil
}
