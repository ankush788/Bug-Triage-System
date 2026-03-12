package service

import (
	"context"
	"errors"
	"time"

	"bug_triage/internal/auth"
	"bug_triage/internal/dto"
	errortype "bug_triage/internal/error"
	"bug_triage/internal/models"
	"bug_triage/internal/repository"
)

// UserService handles user-related business logic
type UserService struct {
	repo             repository.UserRepository
	passwordManager  *auth.PasswordManager
	jwtManager       *auth.JWTManager
	tokenExpireHours int
}

func NewUserService(
	repo repository.UserRepository,
	passwordManager *auth.PasswordManager,
	jwtManager *auth.JWTManager,
) *UserService {
	return &UserService{
		repo:             repo,
		passwordManager:  passwordManager,
		jwtManager:       jwtManager,
		tokenExpireHours: 24, // 24 hour tokens
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

// Login authenticates a user and returns a token
func (s *UserService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	// Get user by email
	user, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, errortype.ErrNotFound) {
			return nil, errors.New("invalid email or password")
		}
		return nil, err
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
