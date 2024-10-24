package services

import (
	"context"
	"errors"

	"taskmaster/internal/models"
	"taskmaster/internal/repositories"

	"golang.org/x/crypto/bcrypt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type UserService interface {
   CreateUser(ctx context.Context, user *models.User) error
   Authenticate(ctx context.Context, username, password string) (*models.User, error)
}

type userService struct {
   userRepo repositories.UserRepository
   logger   zerolog.Logger
}

func NewUserService(repo repositories.UserRepository) UserService {
   serviceLogger := log.With().Str("component", "user_service").Logger()

   return &userService{
      userRepo: repo,
      logger: serviceLogger,
   }
}

func (s *userService) CreateUser(ctx context.Context, user *models.User) error {
   // NOTE: business rules will come here, such as checking user validity
   if user.Username == "" {
      return errors.New("Username is required")
   }
   if user.Email == "" {
      return errors.New("Email is required")
   }
   if user.Password == "" {
      return errors.New("Password is required")
   }

   s.logger.Info().Msg("CreateUser called")

   // Call the repository to create the user and storing it in db
   return s.userRepo.CreateUser(ctx, user)
}

func (s *userService) Authenticate(ctx context.Context, username, password string) (*models.User, error) {
   user, err := s.userRepo.GetUserByUsername(ctx, username)
   if err != nil {
      return nil, err
   }

   if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
      return nil, errors.New("invalid credentials")
   }

   user.Password = ""  // clear the hashed password before returning
   return user, nil
}
