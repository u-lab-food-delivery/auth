package service

import (
	"auth_service/genproto/auth"
	"auth_service/storage/postgres"
	"context"
)

type AuthService struct {
	user *postgres.UserManagementImpl
	auth.UnimplementedAuthServiceServer
}

func (auth.UnimplementedAuthServiceServer) Login(context.Context, *auth.LoginRequest) (*auth.LoginResponse, error)
func (auth.UnimplementedAuthServiceServer) RefreshToken(context.Context, *auth.RefreshTokenRequest) (*auth.RefreshTokenResponse, error)
func (auth.UnimplementedAuthServiceServer) Register(context.Context, *auth.RegisterRequest) (*auth.RegisterResponse, error)
func (auth.UnimplementedAuthServiceServer) VerifyEmail(context.Context, *auth.VerifyEmailRequest) (*auth.VerifyEmailResponse, error)
