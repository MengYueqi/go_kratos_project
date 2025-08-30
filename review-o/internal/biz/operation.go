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
type Operation struct {
	Hello string
}

// 定义 O 端操作数据
type OperaParam struct {
	ID       int64
	AppealID int64
	Reason   string
	Status   int32
	OpUser   string
}

// OperationRepo is a Greater repo.
type OperationRepo interface {
	Save(context.Context, *Operation) (*Operation, error)
	Update(context.Context, *Operation) (*Operation, error)
	FindByID(context.Context, int64) (*Operation, error)
	ListByHello(context.Context, string) ([]*Operation, error)
	ListAll(context.Context) ([]*Operation, error)
	AppealOperate(context.Context, *OperaParam) (*OperaParam, error)
}

// OperationUsecase is a Greeter usecase.
type OperationUsecase struct {
	repo OperationRepo
	log  *log.Helper
}

// NewGreeterUsecase new a Greeter usecase.
func NewOperationUsecase(repo OperationRepo, logger log.Logger) *OperationUsecase {
	return &OperationUsecase{repo: repo, log: log.NewHelper(logger)}
}

// CreateOperation creates a Greeter, and returns the new Greeter.
func (uc *OperationUsecase) CreateOperation(ctx context.Context, g *Operation) (*Operation, error) {
	uc.log.WithContext(ctx).Infof("CreateOperation: %v", g.Hello)
	return uc.repo.Save(ctx, g)
}

func (uc *OperationUsecase) OperateAppeal(ctx context.Context, op *OperaParam) (*OperaParam, error) {
	opAppeal, err := uc.repo.AppealOperate(ctx, op)
	if err != nil {
		return nil, err
	}
	return opAppeal, nil
}
