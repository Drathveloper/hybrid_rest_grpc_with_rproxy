package controller

import (
	"context"
	"net/http"
	"poc_hybrid_grpc_rest/customerror"
	"poc_hybrid_grpc_rest/model"
	"poc_hybrid_grpc_rest/service"
)

type CreateUserRequest struct {
	Username string   `json:"username"`
	Password string   `json:"password"`
	Roles    []string `json:"roles"`
}

type CreateUserResponse struct {
	ID string `json:"id"`
}

type ValidateAccessTokenRequest struct {
	AccessToken string `json:"accessToken"`
}

type ValidateAccessTokenResponse struct {
	IdentityToken string `json:"identityToken"`
}

type UserRestController interface {
	CreateUser(ctx context.Context, request CreateUserRequest) (CreateUserResponse, *customerror.RestError)
	Validate(ctx context.Context, request ValidateAccessTokenRequest) (ValidateAccessTokenResponse, *customerror.RestError)
}

type userRestController struct {
	userService service.UserService
}

func NewUserRestController(service service.UserService) UserRestController {
	return &userRestController{
		userService: service,
	}
}

func (c *userRestController) CreateUser(ctx context.Context, request CreateUserRequest) (CreateUserResponse, *customerror.RestError) {
	headers, ok := ctx.Value(HttpHeaderKey{}).(http.Header)
	if !ok {
		return CreateUserResponse{}, &customerror.RestError{Code: http.StatusInternalServerError, Message: "error while getting headers"}
	}
	identityToken := headers.Get("identityToken")
	if identityToken == "" {
		return CreateUserResponse{}, &customerror.RestError{Code: http.StatusUnauthorized, Message: "identityToken cannot be empty"}
	}
	ID, err := c.userService.CreateUser(model.User{
		Username: request.Username,
		Password: request.Password,
		Roles:    request.Roles,
	})
	if err != nil {
		switch err {
		case service.ErrInvalidUser:
			return CreateUserResponse{}, &customerror.RestError{Code: http.StatusBadRequest, Message: err.Error()}
		default:
			return CreateUserResponse{}, &customerror.RestError{Code: http.StatusInternalServerError, Message: err.Error()}
		}
	}
	return CreateUserResponse{
		ID: identityToken + "|" + ID,
	}, nil
}

func (c *userRestController) Validate(ctx context.Context, request ValidateAccessTokenRequest) (ValidateAccessTokenResponse, *customerror.RestError) {
	if request.AccessToken == "123123" {
		return ValidateAccessTokenResponse{
			IdentityToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
		}, nil
	}
	return ValidateAccessTokenResponse{}, &customerror.RestError{Code: http.StatusUnauthorized, Message: "invalid accessToken"}
}
