package user

import (
	"context"
	"demo-service/helpers"
	"demo-service/proto/pb"
	"demo-service/services/user/business"
	"demo-service/services/user/entity"
	"fmt"
	"github.com/viettranx/service-context/core"
	"log/slog"
	"os"
)

var logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

type Business interface {
	GetUserDetails(ctx context.Context, id int) (*pb.User, error)
	GetUsersByIds(ctx context.Context, ids []int) ([]pb.User, error)
	CreateNewUser(ctx context.Context, data *entity.UserDataCreation) error
}

type grpcService struct {
	business   Business
	repository business.UserRepository
	time       helpers.Timer
}

func (s *grpcService) GetUserProfile(ctx context.Context) (*pb.User, error) {
	requester := core.GetRequester(ctx)
	method := "GetUserProfile"
	s.time.Start()
	logger.Info("request", "method", method)
	if requester == nil {
		return nil, core.ErrForbidden.
			WithError(entity.ErrCannotGetUser.Error())
	}

	uid, _ := core.FromBase58(requester.GetSubject())
	requesterId := int(uid.GetLocalID())

	user, err := s.repository.GetUserById(ctx, requesterId)

	if err != nil {
		logger.Error("response", "method", method, "err", err, "ms", s.time.End())
		return nil, core.ErrUnauthorized.
			WithError(entity.ErrCannotGetUser.Error()).
			WithDebug(err.Error())
	}
	logger.Info("response", "method", method, "data", user, "ms", s.time.End())
	return user, nil
}

func NewService(business Business, repository business.UserRepository) *grpcService {
	return &grpcService{business: business, repository: repository}
}

func (s *grpcService) GetUserById(ctx context.Context, req *pb.GetUserByIdReq) (*pb.PublicUserInfoResp, error) {
	method := "GetUserById"
	s.time.Start()
	logger.Info("request", "method", method)
	user, err := s.business.GetUserDetails(ctx, int(req.Id))

	if err != nil {
		logger.Error("response", "method", method, "err", err, "ms", s.time.End())
		return nil, core.ErrInternalServerError.WithError(err.Error())
	}
	logger.Info("response", "method", method, "data", user, "ms", s.time.End())
	return &pb.PublicUserInfoResp{
		User: &pb.PublicUserInfo{
			Id:        int32(user.Id),
			FirstName: user.FirstName,
			LastName:  user.LastName,
		},
	}, nil
}

func (s *grpcService) GetUserDetailsById(ctx context.Context, req *pb.GetUserByIdReq) (*pb.PublicUserInfoResp, error) {
	method := "GetUserDetailsById"
	s.time.Start()
	logger.Info("request", "method", method)
	user, err := s.business.GetUserDetails(ctx, int(req.Id))

	if err != nil {
		logger.Error("response", "method", method, "err", err, "ms", s.time.End())
		return nil, core.ErrInternalServerError.WithError(err.Error())
	}
	logger.Info("response", "method", method, "data", user, "ms", s.time.End())
	return &pb.PublicUserInfoResp{
		User: &pb.PublicUserInfo{
			Id:        int32(user.Id),
			FirstName: user.FirstName,
			LastName:  user.LastName,
		},
	}, nil
}

func (s *grpcService) GetUsersByIds(ctx context.Context, req *pb.GetUsersByIdsReq) (*pb.PublicUsersInfoResp, error) {
	method := "GetUsersByIds"
	s.time.Start()
	logger.Info("request", "method", method)
	userIDs := make([]int, len(req.Ids))

	for i := range userIDs {
		userIDs[i] = int(req.Ids[i])
	}

	fmt.Println("userIDs", userIDs)

	users, err := s.business.GetUsersByIds(ctx, userIDs)

	if err != nil {
		logger.Error("response", "method", method, "err", err, "ms", s.time.End())
		return nil, core.ErrInternalServerError.WithError(err.Error())
	}

	publicUserInfo := make([]*pb.PublicUserInfo, len(users))

	for i := range users {
		publicUserInfo[i] = &pb.PublicUserInfo{
			Id:        int32(users[i].Id),
			FirstName: users[i].FirstName,
			LastName:  users[i].LastName,
		}
	}
	logger.Info("response", "method", method, "data", publicUserInfo, "ms", s.time.End())
	return &pb.PublicUsersInfoResp{Users: publicUserInfo}, nil
}

func (s *grpcService) CreateUser(ctx context.Context, req *pb.CreateUserReq) (*pb.NewUserIdResp, error) {
	method := "CreateUser"
	s.time.Start()
	logger.Info("request", "method", method)
	newUserData := entity.NewUserForCreation(req.FirstName, req.LastName, req.Email)

	if err := s.business.CreateNewUser(ctx, &newUserData); err != nil {
		logger.Error("response", "method", method, "err", err, "ms", s.time.End())
		return nil, core.ErrInternalServerError.WithError(err.Error())
	}
	logger.Info("response", "method", method, "data", newUserData, "ms", s.time.End())
	return &pb.NewUserIdResp{Id: int32(newUserData.Id)}, nil
}
