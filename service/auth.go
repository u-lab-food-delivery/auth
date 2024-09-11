package service

import (
	"auth_service/config"
	"auth_service/genproto/auth"
	"auth_service/models"
	"auth_service/storage/cache"
	"auth_service/storage/postgres"
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type AuthService struct {
	user        *postgres.UserManagementImpl
	emailsender *EmailSender
	cnf         *config.Config
	tokenCache  *cache.TokenCache
	auth.UnimplementedAuthServiceServer
}

func NewAuthService(user *postgres.UserManagementImpl, emailsender *EmailSender, cnf *config.Config, tokenCacher *cache.TokenCache) *AuthService {
	return &AuthService{
		user:        user,
		emailsender: emailsender,
		cnf:         cnf,
		tokenCache:  tokenCacher,
	}
}

func (a *AuthService) CheckByEmail(ctx context.Context, req *auth.CheckByEmailRequest) (*auth.EmptyMessage, error) {
	user, err := a.user.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	if user != nil {
		return &auth.EmptyMessage{}, nil
	}

	return nil, nil
}

func (a *AuthService) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	user, err := a.user.GetByEmail(ctx, req.Email)
	if err != nil {
		log.Println("Unexpected error has occured: ", err)
		return nil, err
	}
	if err == sql.ErrNoRows && user == nil {
		log.Println("No user foun by this email: ", req.Email)
		return nil, fmt.Errorf("notFound")
	}

	tokens, err := a.CreateToken(ctx, &auth.CreateTokenRequest{UserId: user.UserId})

	if err != nil {
		return nil, err
	}

	return &auth.LoginResponse{AccessToken: tokens.AccessToken, RefreshToken: tokens.RefreshToken}, nil
}

func (a *AuthService) LogOut(ctx context.Context, req *auth.LogOutRequest) (*auth.EmptyMessage, error) {
	err := a.tokenCache.RevokeToken(ctx, &cache.RevokeTokens{RefreshToken: req.RefreshToken})
	if err != nil {
		log.Println("Token Revocation failed: ", err)
		return nil, err
	}

	return &auth.EmptyMessage{}, nil
}

func (a *AuthService) RefreshToken(ctx context.Context, req *auth.RefreshTokenRequest) (*auth.RefreshTokenResponse, error) {
	ok, err := a.tokenCache.IsRevokedToken(ctx, req.RefreshToken, req.RefreshToken)
	if err != nil {
		log.Println("Checking token faild: ", err)
		return nil, err
	}

	if !ok {
		return nil, fmt.Errorf("Refresh token revoked")
	}

	claims, err := cache.ExtractClaims(req.RefreshToken)
	if err != nil {
		log.Println("Couldn't extract claims: ", err)
		return nil, err
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, models.Claims{
		UserID: claims.UserID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 5)),
		},
	})

	accessTokenStr, err := accessToken.SignedString([]byte(a.cnf.JWT.SecretKey))
	if err != nil {
		log.Println("Failed to create access token: ", err)
		return nil, err
	}

	return &auth.RefreshTokenResponse{AccessToken: accessTokenStr}, nil
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

func (a *AuthService) RevokeToken(ctx context.Context, req *auth.RevokeTokenRequest) (*auth.RevokeTokenResponse, error) {
	if err := a.tokenCache.RevokeToken(ctx, &cache.RevokeTokens{AccessToken: req.Token}); err != nil {
		log.Println("Revocation access token failed: ", err)
		return nil, err
	}

	return &auth.RevokeTokenResponse{Message: "Access token revoked"}, nil
}

func (a *AuthService) CreateToken(ctx context.Context, req *auth.CreateTokenRequest) (*auth.CreateTokenResponse, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, models.Claims{
		UserID: req.UserId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 5)),
		},
	})

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, models.Claims{
		UserID: req.UserId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	})

	accessTokenStr, err := accessToken.SignedString([]byte(a.cnf.JWT.SecretKey))
	if err != nil {
		log.Println("Failed to create access token: ", err)
		return nil, err
	}
	RefreshTokenStr, err := refreshToken.SignedString([]byte(a.cnf.JWT.SecretKey))
	if err != nil {
		log.Println("Failed to create refresh token: ", err)
		return nil, err
	}

	return &auth.CreateTokenResponse{
		AccessToken:  accessTokenStr,
		RefreshToken: RefreshTokenStr,
	}, nil
}

func (a *AuthService) GetToken(ctx context.Context, req *auth.GetTokenRequest) (*auth.GetTokenResponse, error) {

	return nil, nil
}
