package service

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"go-admin/server/internal/config"
	"go-admin/server/internal/model"
	"go-admin/server/internal/repository"
	"go-admin/server/internal/utils"
	"go-admin/server/pkg/jwtx"
)

type AuthService struct {
	repo   *repository.Repository
	redis  *redis.Client
	config *config.Config
}

func NewAuthService(repo *repository.Repository, redis *redis.Client, cfg *config.Config) *AuthService {
	return &AuthService{repo: repo, redis: redis, config: cfg}
}

type LoginPayload struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type TokenPair struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    int    `json:"expiresIn"`
}

func (s *AuthService) Login(ctx context.Context, payload LoginPayload) (*model.User, *TokenPair, error) {
	user, err := s.repo.GetUserByUsername(ctx, payload.Username)
	if err != nil {
		return nil, nil, err
	}
	if user.Status != 1 {
		return nil, nil, fmt.Errorf("user disabled")
	}
	if err := utils.CheckPassword(user.Password, payload.Password); err != nil {
		return nil, nil, fmt.Errorf("invalid credentials")
	}
	tokens, err := s.issueTokenPair(ctx, user)
	return user, tokens, err
}

func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (*model.User, *TokenPair, error) {
	claims, err := jwtx.Parse(s.config.JWTSecret, refreshToken)
	if err != nil {
		return nil, nil, err
	}
	key := s.refreshTokenKey(claims.UserID)
	storedToken, err := s.redis.Get(ctx, key).Result()
	if err != nil {
		return nil, nil, err
	}
	if storedToken != refreshToken {
		return nil, nil, fmt.Errorf("refresh token expired")
	}
	user, err := s.repo.GetUserByID(ctx, claims.UserID)
	if err != nil {
		return nil, nil, err
	}
	if user.Status != 1 {
		return nil, nil, fmt.Errorf("user disabled")
	}
	tokens, err := s.issueTokenPair(ctx, user)
	return user, tokens, err
}

func (s *AuthService) Logout(ctx context.Context, userID uint) error {
	return s.redis.Del(ctx, s.refreshTokenKey(userID)).Err()
}

func (s *AuthService) issueTokenPair(ctx context.Context, user *model.User) (*TokenPair, error) {
	accessExpire := time.Duration(s.config.JWTAccessExpireMinute) * time.Minute
	refreshExpire := time.Duration(s.config.JWTRefreshExpireHour) * time.Hour

	accessToken, err := jwtx.Generate(s.config.JWTSecret, user.ID, user.Username, accessExpire)
	if err != nil {
		return nil, err
	}
	refreshToken, err := jwtx.Generate(s.config.JWTSecret, user.ID, user.Username, refreshExpire)
	if err != nil {
		return nil, err
	}
	if err := s.redis.Set(ctx, s.refreshTokenKey(user.ID), refreshToken, refreshExpire).Err(); err != nil {
		return nil, err
	}
	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int(accessExpire.Seconds()),
	}, nil
}

func (s *AuthService) refreshTokenKey(userID uint) string {
	return fmt.Sprintf("refresh_token:%d", userID)
}
