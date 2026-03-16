package service

import (
	"context"
	"errors"
	"time"

	"bug_triage/internal/auth"
	"bug_triage/internal/cache"
	"bug_triage/internal/dto"
	errortype "bug_triage/internal/error"
	"bug_triage/internal/models"
	"bug_triage/internal/repository"

	"go.uber.org/zap"
)

// UserService handles user-related business logic
type UserService struct {
	repo             repository.UserRepository
	passwordManager  *auth.PasswordManager
	jwtManager       *auth.JWTManager
	tokenExpireHours int
	userCache        *cache.UserCache
	logger           *zap.Logger
}

func NewUserService(
	repo repository.UserRepository,
	passwordManager *auth.PasswordManager,
	jwtManager *auth.JWTManager,
	userCache *cache.UserCache,
	logger *zap.Logger,
) *UserService {
	return &UserService{
		repo:             repo,
		passwordManager:  passwordManager,
		jwtManager:       jwtManager,
		tokenExpireHours: 24, // 24 hour tokens
		userCache:        userCache,
		logger:           logger,
	}
}

// Register creates a new user account
func (s *UserService) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.RegisterResponse, error) {
	// Check if user already exists
	existing, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, errortype.ErrNotFound) {
			// no existing user, continue
		} else {
			return nil, err
		}
	}
	if existing != nil {
		return nil, errors.New("user already exists")
	}

	// Hash password
	hashedPassword, err := s.passwordManager.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &models.User{
		Email:        req.Email,
		PasswordHash: hashedPassword,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Cache the user
	if err := s.userCache.Set(ctx, user.Email, user); err != nil {
		s.logger.Warn("failed to cache user", zap.Error(err))
	}

	// Generate JWT token
	token, err := s.jwtManager.GenerateToken(user.ID, user.Email, time.Duration(s.tokenExpireHours)*time.Hour)
	if err != nil {
		return nil, err
	}

	return &dto.RegisterResponse{
		ID:    user.ID,
		Email: user.Email,
		Token: token,
	}, nil
}

// GetUser retrieves a user by ID
func (s *UserService) GetUser(ctx context.Context, userID int64) (*models.User, error) {
	// Get from database
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, errortype.ErrNotFound) {
			return nil, errortype.ErrNotFound
		}
		s.logger.Error("failed to get user", zap.Error(err))
		return nil, err
	}

	// Cache the result
	if err := s.userCache.Set(ctx, user.Email, user); err != nil {
		s.logger.Warn("failed to cache user", zap.Error(err))
	}

	return user, nil
}

// Login authenticates a user and returns a token
func (s *UserService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	// Try to get from cache
	user, err := s.userCache.Get(ctx, req.Email)
	if err == nil {
		s.logger.Debug("user retrieved from cache", zap.String("email", req.Email))
		// Verify password
		if !s.passwordManager.VerifyPassword(user.PasswordHash, req.Password) {
			return nil, errors.New("invalid email or password")
		}

		// Generate JWT token
		token, err := s.jwtManager.GenerateToken(user.ID, user.Email, time.Duration(s.tokenExpireHours)*time.Hour)
		if err != nil {
			return nil, err
		}

		return &dto.LoginResponse{
			ID:    user.ID,
			Email: user.Email,
			Token: token,
		}, nil
	}

	// Get user by email
	user, err = s.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, errortype.ErrNotFound) {
			return nil, errors.New("invalid email or password")
		}
		return nil, err
	}

	// Cache the user
	if err := s.userCache.Set(ctx, user.Email, user); err != nil {
		s.logger.Warn("failed to cache user", zap.Error(err))
	}

	// Verify password
	if !s.passwordManager.VerifyPassword(user.PasswordHash, req.Password) {
		return nil, errors.New("invalid email or password")
	}

	// Generate JWT token
	token, err := s.jwtManager.GenerateToken(user.ID, user.Email, time.Duration(s.tokenExpireHours)*time.Hour)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		ID:    user.ID,
		Email: user.Email,
		Token: token,
	}, nil
}
