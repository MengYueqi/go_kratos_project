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

// 创建一条评论
func (r *ReviewerRepo) AddReviewReply(ctx context.Context, reply *model.ReviewReplyInfo) (int64, error) {
	// 查询 ShoreID 是否与评论的 Review 中的一致
	rv, err := r.data.query.ReviewInfo.
		WithContext(ctx).
		Where(r.data.query.ReviewInfo.ReviewID.Eq(reply.ReviewID)).
		Find()
	if err != nil {
		return 0, v1.ErrorDbFailed("DB error while searching reviewID: %v", reply.ReviewID)
	}
	if len(rv) == 0 {
		return 0, v1.ErrorReviewidErr("Do not exist ReviewID: %v", reply.ReviewID)
	}
	// 处理 StoreID 和评论不一致的情况
	if rv[0].StoreID != reply.StoreID {
		return 0, v1.ErrorStoreidReviewidMismatch("Store ID mismatch with View's, StoreID: %v, View's StoreID: %v", reply.StoreID, rv[0].StoreID)
	}
	// 核心逻辑
	err = r.data.query.ReviewReplyInfo.WithContext(ctx).Save(reply)
	if err != nil {
		return 0, v1.ErrorDbFailed("DB Save error")
	}
	return reply.ReplyID, nil
}

func (r *ReviewerRepo) AddAppealReview(ctx context.Context, appeal *model.ReviewAppealInfo) (int64, error) {
	// 插入一条申诉记录
	err := r.data.query.ReviewAppealInfo.
		WithContext(ctx).
		Save(appeal)
	if err != nil {
		return 0, v1.ErrorDbFailed("DB Save error")
	}
	return appeal.AppealID, nil
}

func (r *ReviewerRepo) GetAppealByReviewID(ctx context.Context, reviewID int64) ([]*model.ReviewAppealInfo, error) {
	data, err := r.data.query.ReviewAppealInfo.
		WithContext(ctx).
		Where(r.data.query.ReviewAppealInfo.ReviewID.Eq(reviewID)).
		Find()
	if err != nil {
		return nil, v1.ErrorDbFailed("DB error while finding %v", reviewID)
	}
	return data, nil
}

func (r *ReviewerRepo) UpdateAppealByAppealID(ctx context.Context, appeal *model.ReviewAppealInfo) (*model.ReviewAppealInfo, error) {
	_, err := r.data.query.ReviewAppealInfo.
		WithContext(ctx).
		Where(r.data.query.ReviewAppealInfo.AppealID.Eq(appeal.AppealID)).
		Updates(appeal)
	if err != nil {
		return &model.ReviewAppealInfo{}, err
	}
	return appeal, nil
}

// 通过申诉 ID 获取申诉信息
func (r *ReviewerRepo) GetAppealByAppealID(ctx context.Context, appealID int64) ([]*model.ReviewAppealInfo, error) {
	info, err := r.data.query.ReviewAppealInfo.
		WithContext(ctx).
		Where(r.data.query.ReviewAppealInfo.AppealID.Eq(appealID)).
		Find()
	if err != nil {
		return nil, v1.ErrorIdErr("Do not exist AppealID: %v", appealID)
	}
	return info, nil
}
