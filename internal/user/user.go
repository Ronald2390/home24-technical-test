package user

import (
	"context"
	"errors"
	"time"

	"home24-technical-test/internal/user/model"
	"home24-technical-test/internal/user/public"
	"home24-technical-test/pkg/appcontext"

	"golang.org/x/crypto/bcrypt"
)

// Storage represents the user storage interface
type Storage interface {
	FindAll(ctx context.Context, params *public.FindAllUsersParams) ([]*model.User, error)
	FindByID(ctx context.Context, userID int) (*model.User, error)
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	Insert(ctx context.Context, user *model.User) error
	Update(ctx context.Context, updatedUser *model.User) error
	Delete(ctx context.Context, userID int) error
}

// ServiceInterface represents the user service interface
type ServiceInterface interface {
	ListUsers(ctx context.Context, params *public.FindAllUsersParams) ([]*model.User, error)
	GetUser(ctx context.Context, userID int) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	CreateUser(ctx context.Context, params *public.CreateUserParams) error
	UpdateUser(ctx context.Context, params *public.UpdateUserParams) (*model.User, error)
	DeleteUser(ctx context.Context, userID int) error
	ChangePassword(ctx context.Context, userID int, oldPassword, newPassword string) error
}

// Errors
var (
	ErrWrongPassword      = errors.New("wrong password")
	ErrWrongEmail         = errors.New("wrong email")
	ErrEmailAlreadyExists = errors.New("Email Already Exists")
	ErrNotFound           = errors.New("not found")
	ErrNoInput            = errors.New("no input")
)

// Service is the domain logic implementation of user Service interface
type Service struct {
	repository Storage
}

// ListUsers is listing all Users
func (s *Service) ListUsers(ctx context.Context, params *public.FindAllUsersParams) ([]*model.User, error) {
	users, err := s.repository.FindAll(ctx, params)
	if err != nil {
		return nil, err
	}

	return users, nil
}

// GetUser get user by its id
func (s *Service) GetUser(ctx context.Context, userID int) (*model.User, error) {
	user, err := s.repository.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByEmail get user by its email
func (s *Service) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	user, err := s.repository.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// UpdateUser updates users data
func (s *Service) UpdateUser(ctx context.Context, params *public.UpdateUserParams) (*model.User, error) {
	user, err := s.repository.FindByEmail(ctx, params.Email)
	if err != nil {
		return nil, err
	}

	if user != nil && user.ID != params.ID {
		return nil, ErrEmailAlreadyExists
	}

	updatedUser, err := s.repository.FindByID(ctx, params.ID)
	if err != nil {
		return nil, err
	}

	if params.Name != "" {
		updatedUser.Name = params.Name
	}
	if params.Email != "" {
		updatedUser.Email = params.Email
	}
	if params.Address != "" {
		updatedUser.Address = params.Address
	}

	currentUserID := appcontext.UserID(ctx)

	updatedUser.UpdatedBy = currentUserID
	updatedUser.UpdatedAt = time.Now()

	err = s.repository.Update(ctx, updatedUser)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}

// ChangePassword changes user's password
func (s *Service) ChangePassword(ctx context.Context, userID int, oldPassword, newPassword string) error {
	currentUser, err := s.repository.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(currentUser.Password), []byte(oldPassword))
	if err != nil {
		return ErrWrongPassword
	}

	bcryptHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	currentUser.Password = string(bcryptHash)
	err = s.repository.Update(ctx, currentUser)
	if err != nil {
		return err
	}

	return nil
}

//DeleteUser deleting user and its session
func (s *Service) DeleteUser(ctx context.Context, userID int) error {
	singleUser, err := s.repository.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	err = s.repository.Delete(ctx, singleUser.ID)
	if err != nil {
		return err
	}

	return nil
}

// CreateUser creates a new user
func (s *Service) CreateUser(ctx context.Context, params *public.CreateUserParams) error {
	bcryptHash, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	existingUser, err := s.repository.FindByEmail(ctx, params.Email)
	if err != nil {
		return err
	}

	if existingUser != nil {
		return ErrEmailAlreadyExists
	}

	currentUserID := appcontext.UserID(ctx)

	user := &model.User{
		Name:      params.Name,
		Email:     params.Email,
		Address:   params.Address,
		Password:  string(bcryptHash),
		CreatedBy: currentUserID,
		CreatedAt: time.Now(),
		UpdatedBy: currentUserID,
		UpdatedAt: time.Now(),
	}

	err = s.repository.Insert(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

// NewService creates a new user AppService
func NewService(
	userRepository Storage,
) *Service {
	return &Service{
		repository: userRepository,
	}
}
