package service

import (
	"context"
	"errors"
	"time"

	"bug_triage/internal/auth"
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

// RegisterRequest holds incoming registration data
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// RegisterResponse returns user info and auth token
type RegisterResponse struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
	Token string `json:"token"`
}

// Register creates a new user account
func (s *UserService) Register(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error) {
	// Check if user already exists
	existing, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
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

	return &RegisterResponse{
		ID:    user.ID,
		Email: user.Email,
		Token: token,
	}, nil
}

// LoginRequest holds incoming login data
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse returns user info and auth token
type LoginResponse struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
	Token string `json:"token"`
}

// Login authenticates a user and returns a token
func (s *UserService) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	// Get user by email
	user, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid email or password")
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

	return &LoginResponse{
		ID:    user.ID,
		Email: user.Email,
		Token: token,
	}, nil
}
