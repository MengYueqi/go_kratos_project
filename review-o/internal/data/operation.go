package data

import (
	"context"
	reviewv1 "review-o/api/review/v1"

	"review-o/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type operateRepo struct {
	data *Data
	log  *log.Helper
}

// NewOperationRepo .
func NewOperationRepo(data *Data, logger log.Logger) biz.OperationRepo {
	return &operateRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *operateRepo) Save(ctx context.Context, g *biz.Operation) (*biz.Operation, error) {
	return g, nil
}

func (r *operateRepo) Update(ctx context.Context, g *biz.Operation) (*biz.Operation, error) {
	return g, nil
}

func (r *operateRepo) FindByID(context.Context, int64) (*biz.Operation, error) {
	return nil, nil
}

func (r *operateRepo) ListByHello(context.Context, string) ([]*biz.Operation, error) {
	return nil, nil
}

func (r *operateRepo) ListAll(context.Context) ([]*biz.Operation, error) {
	return nil, nil
}

func (r *operateRepo) AppealOperate(ctx context.Context, op *biz.OperaParam) (*biz.OperaParam, error) {
	appeal, err := r.data.rc.HandleAppeal(context.Background(), &reviewv1.AppealOperateRequest{
		AppealID: op.AppealID,
		Reason:   op.Reason,
		Status:   op.Status,
		OpUser:   op.OpUser,
		ID:       op.ID,
	})
	if err != nil {
		return nil, err
	}
	return &biz.OperaParam{
		ID:       appeal.ID,
		AppealID: appeal.AppealID,
		Reason:   appeal.Reason,
		Status:   appeal.Status,
		OpUser:   appeal.OpUser,
	}, nil
}
