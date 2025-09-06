package data

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/go-redis/redis/v8"
	"golang.org/x/sync/singleflight"
	v1 "review-service/api/review/v1"
	"review-service/internal/data/model"
	"strconv"
	"strings"
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

// 根据 StoreID offset 和 limit 获取评论列表
func (r *ReviewerRepo) ListReviewByStoreID(ctx context.Context, storeID int64, offset int32, limit int32) ([]*biz.MyReviewInfo, error) {
	return r.getData(ctx, storeID, offset, limit)
	//// 去 ES 中查询评价
	//resp, err := r.data.es.Search().
	//	Index("review").
	//	From(int(offset)).
	//	Size(int(limit)).Query(&types.Query{
	//	Bool: &types.BoolQuery{
	//		Filter: []types.Query{
	//			{
	//				Term: map[string]types.TermQuery{
	//					"store_id": {Value: storeID},
	//				},
	//			},
	//		},
	//	},
	//}).Do(ctx)
	//if err != nil {
	//	return nil, v1.ErrorDbFailed("ES search error")
	//}
	//
	//rv := make([]*biz.MyReviewInfo, 0, resp.Hits.Total.Value)
	//// 反序列化数据
	//for _, hit := range resp.Hits.Hits {
	//	temp := &biz.MyReviewInfo{}
	//	err := json.Unmarshal(hit.Source_, temp)
	//	if err != nil {
	//		r.log.Errorf("json unmarshal error: %v", err)
	//		continue
	//	}
	//	rv = append(rv, temp)
	//}
	//return rv, nil
}

var g singleflight.Group

func (r *ReviewerRepo) getData(ctx context.Context, storeID int64, offset int32, limit int32) ([]*biz.MyReviewInfo, error) {
	// 使用 singleflight 取数据，防止缓存击穿
	data, err := r.getDataFromSingleFlight(ctx, "review:"+strconv.FormatInt(storeID, 10)+":"+strconv.Itoa(int(offset))+":"+strconv.Itoa(int(limit)))
	if err != nil {
		return nil, err
	}
	hm := new(types.HitsMetadata)
	err = json.Unmarshal(data, hm)
	if err != nil {
		return nil, err
	}
	rv := make([]*biz.MyReviewInfo, 0, hm.Total.Value)

	// 反序列化数据
	for _, hit := range hm.Hits {
		temp := &biz.MyReviewInfo{}
		err := json.Unmarshal(hit.Source_, temp)
		if err != nil {
			r.log.Errorf("json unmarshal error: %v", err)
			continue
		}
		rv = append(rv, temp)
	}
	return rv, nil
}

// 使用 singleflight 防止缓存击穿
// key: review:storeID:page:pageSize
func (r *ReviewerRepo) getDataFromSingleFlight(ctx context.Context, key string) ([]byte, error) {
	v, err, _ := g.Do(key, func() (interface{}, error) {
		// 查询数据库
		data, err := r.getDataFromCache(ctx, key)
		if err == nil {
			r.log.Debugf("getDataFromCache key: %s", key)
			return data, nil
		}
		if errors.Is(err, redis.Nil) {
			// 缓存中没有此数据，需要查询 ES
			esData, err := r.getDataES(ctx, key)
			if err != nil {
				return nil, err
			}
			return esData, r.setCache(ctx, key, esData)
		}
		// redis 查询出错
		return nil, err
	})
	r.log.Debugf("getDataFromSingleFlight key: %s, value: %v, err: %v", key, v, err)
	if err != nil {
		return nil, err
	}
	return v.([]byte), nil
}

func (r *ReviewerRepo) getDataFromCache(ctx context.Context, key string) ([]byte, error) {
	return r.data.redis.Get(ctx, key).Bytes()
}

func (r *ReviewerRepo) setCache(ctx context.Context, key string, data []byte) error {
	return r.data.redis.Set(ctx, key, data, time.Minute*5).Err()
}

func (r *ReviewerRepo) getDataES(ctx context.Context, key string) ([]byte, error) {
	value := strings.Split(key, ":")
	// 对 key 的长度进行检查
	if len(strings.Split(key, ":")) != 4 {
		return nil, errors.New("key format error")
	}
	index, storeIDStr, offsetStr, limitStr := value[0], value[1], value[2], value[3]
	// 进行类型转换
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		return nil, errors.New("offset format error")
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return nil, errors.New("limit format error")
	}
	// 去 ES 中查询评价
	resp, err := r.data.es.Search().
		Index(index).
		From(int(offset)).
		Size(int(limit)).Query(&types.Query{
		Bool: &types.BoolQuery{
			Filter: []types.Query{
				{
					Term: map[string]types.TermQuery{
						"store_id": {Value: storeIDStr},
					},
				},
			},
		},
	}).Do(ctx)
	if err != nil {
		return nil, v1.ErrorDbFailed("ES search error")
	}

	return json.Marshal(resp.Hits)
}
