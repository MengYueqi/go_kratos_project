package data

import (
	"context"
	v1 "review-service/api/review/v1"
	"review-service/internal/data/model"
	"time"

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
	find, err := r.data.query.ReviewInfo.WithContext(ctx).Where(r.data.query.ReviewInfo.OrderID.Eq(orderId)).Find()
	if err != nil {
		return nil, v1.ErrorDbFailed("DB Find error")
	}
	return find, nil
}

func (r *ReviewerRepo) ListByHello(context.Context, string) ([]*biz.Reviewer, error) {
	return nil, nil
}

func (r *ReviewerRepo) ListAll(context.Context) ([]*biz.Reviewer, error) {
	return nil, nil
}

func (r *ReviewerRepo) DeleteReview(ctx context.Context, ID int64) error {
	_, err := r.data.query.ReviewInfo.
		WithContext(ctx).
		Where(r.data.query.ReviewInfo.ID.Eq(ID)).
		Update(r.data.query.ReviewInfo.DeleteAt, time.Now())
	return err
}

func (r *ReviewerRepo) GetReviewByID(ctx context.Context, ID int64) (*model.ReviewInfo, error) {
	info, err := r.data.query.ReviewInfo.
		WithContext(ctx).
		Where(r.data.query.ReviewInfo.ID.Eq(ID)).
		First()
	if err != nil {
		return nil, v1.ErrorIdErr("Do not exist ID: %v", ID)
	}
	return info, nil
}

func (r *ReviewerRepo) GetReviewByReviewID(ctx context.Context, reviewId int64) ([]*model.ReviewInfo, error) {
	info, err := r.data.query.ReviewInfo.
		WithContext(ctx).
		Where(r.data.query.ReviewInfo.ReviewID.Eq(reviewId)).
		Find()
	if err != nil {
		return nil, v1.ErrorDbFailed("DB error while searching reviewID: %v", reviewId)
	}
	return info, nil
}

func (r *ReviewerRepo) UpdateReviewByReviewID(ctx context.Context, rv *model.ReviewInfo) (int64, error) {
	updateReviewData := model.ReviewInfo{
		Content:      rv.Content,
		Score:        rv.Score,
		ServiceScore: rv.ServiceScore,
		ExpressScore: rv.ExpressScore,
		PicInfo:      rv.PicInfo,
		VideoInfo:    rv.VideoInfo,
		Anonymous:    rv.Anonymous,
	}
	_, err := r.data.query.ReviewInfo.
		WithContext(ctx).
		Where(r.data.query.ReviewInfo.ReviewID.Eq(rv.ReviewID)).
		Updates(updateReviewData)
	if err != nil {
		return 0, v1.ErrorIdErr("Do not exist reviewed: %v", rv.ReviewID)
	}
	return rv.ReviewID, nil

}

func (r *ReviewerRepo) GetReviewByUID(ctx context.Context, uid int64) ([]*model.ReviewInfo, error) {
	data, err := r.data.query.ReviewInfo.
		WithContext(ctx).
		Where(r.data.query.ReviewInfo.UserID.Eq(uid)).
		Find()
	if err != nil {
		return nil, v1.ErrorIdErr("DB error while finding %v", uid)
	}
	return data, nil
}
