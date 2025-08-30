package service

import (
	"context"
	v1 "review-o/api/operation/v1"
	"review-o/internal/biz"
)

// OperationService is a greeter service.
type OperationService struct {
	v1.UnimplementedOperationServer

	uc *biz.OperationUsecase
}

// NewGreeterService new a greeter service.
func NewOperationService(uc *biz.OperationUsecase) *OperationService {
	return &OperationService{uc: uc}
}

func (s *OperationService) OperateAppeal(ctx context.Context, req *v1.AppealOperateUserRequest) (*v1.AppealOperateUserReply, error) {
	data, err := s.uc.OperateAppeal(ctx, &biz.OperaParam{
		ID:       req.ID,
		AppealID: req.AppealID,
		Reason:   req.Reason,
		Status:   req.Status,
		OpUser:   req.OpUser,
	})
	if err != nil {
		return nil, err
	}
	return &v1.AppealOperateUserReply{
		AppealID: data.AppealID,
		Status:   data.Status,
		ID:       data.ID,
		Reason:   data.Reason,
		OpUser:   data.OpUser,
	}, nil

}
