package controller

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"poc_hybrid_grpc_rest/model"
	"poc_hybrid_grpc_rest/pb"
	"poc_hybrid_grpc_rest/service"
)

type userGrpcController struct {
	pb.UnimplementedUserServiceServer
	userService service.UserService
}

func NewUserGrpcController(service service.UserService) pb.UserServiceServer {
	return &userGrpcController{
		userService: service,
	}
}

func (c *userGrpcController) CreateUser(ctx context.Context, request *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	identityToken := ""
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		values := md.Get("identityToken")
		if len(values) == 0 {
			return nil, status.New(codes.Unauthenticated, "invalid identity token").Err()
		}
		identityToken = values[0]
	}
	roles := make([]string, 0)
	for _, r := range request.Roles {
		roles = append(roles, r.String())
	}
	ID, err := c.userService.CreateUser(model.User{
		Username: request.Username,
		Password: request.Password,
		Roles:    roles,
	})
	if err != nil {
		return nil, status.New(codes.Unknown, err.Error()).Err()
	}
	return &pb.CreateUserResponse{
		Id: identityToken + "|" + ID,
	}, nil
}
