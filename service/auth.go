package service

import (
	"auth_service/genproto/auth"
	"auth_service/models"
	"auth_service/storage/postgres"
	"context"
	"log"

	"github.com/google/uuid"
)

type AuthService struct {
	user        *postgres.UserManagementImpl
	emailsender *EmailSender
	auth.UnimplementedAuthServiceServer
}

func NewAuthService(user *postgres.UserManagementImpl, emailsender *EmailSender) *AuthService {
	return &AuthService{
		user:        user,
		emailsender: emailsender,
	}
}

func (a *AuthService) CheckByEmail(ctx context.Context, req *auth.CheckByEmailRequest) (*auth.Void, error) {
	user, err := a.user.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	if user != nil {
		return &auth.Void{}, nil
	}

	return nil, nil
}
func (a *AuthService) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {

	return nil, nil
}

func (a *AuthService) RefreshToken(ctx context.Context, req *auth.RefreshTokenRequest) (*auth.RefreshTokenResponse, error) {
	return nil, nil
}

func (a *AuthService) Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	newID := uuid.NewString()
	user, err := a.user.CreateUser(ctx, &models.User{
		UserId:         newID,
		Email:          req.Email,
		HashedPassword: req.Password,
		IsVerified:     false,
		Name:           "user",
	})

	if err != nil {
		log.Println("failed to create user: ", err)
		return nil, err
	}

	err = a.emailsender.SendVerificationEmail(req.Email, req.VerificationLink)
	if err != nil {
		return nil, err
	}

	return &auth.RegisterResponse{
		UserId: newID,
		Email:  user.Email,
	}, nil
}

func (a *AuthService) VerifyEmail(ctx context.Context, req *auth.VerifyEmailRequest) (*auth.VerifyEmailResponse, error) {
	email, err := a.emailsender.cache.GetEmailByLink(req.Token)
	if err != nil {
		return nil, err
	}

	if email == "" {
		return &auth.VerifyEmailResponse{Message: "The link expired"}, nil
	}

	err = a.user.VerifiyUser(ctx, email)
	if err != nil {
		return &auth.VerifyEmailResponse{Message: "Couldn't verified user. Unexpected error occured!"}, err
	}

	return &auth.VerifyEmailResponse{Message: "Verification successfull"}, nil
}
