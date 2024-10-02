package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"taskmaster/internal/models"

   "golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
   CreateUser(ctx context.Context, user *models.User) error
   GetUserByID(ctx context.Context, id int64) (*models.User, error)
   GetUserByUsername(ctx context.Context, username string) (*models.User, error)
   AuthenticateUser(ctx context.Context, username, password string) (*models.User, error)
}

type userRepository struct {
   db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
   return &userRepository{db: db}
}

func (r *userRepository) CreateUser(ctx context.Context, user *models.User) error {
   // Hash the password
   hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
   if err != nil {
      return fmt.Errorf("failed to hash the password: %w", err)
   }

   query := `INSERT INTO users (username, email, password) VALUES (?, ?, ?)`
   result, err := r.db.ExecContext(ctx, query, user.Username, user.Email, string(hashedPassword))
   if err != nil {
      return fmt.Errorf("failed to create user: %w", err)
   }

   id, err := result.LastInsertId()
   if err != nil {
      return fmt.Errorf("failed to get last insert id: %w", err)
   }

   user.ID = id
   user.Password = ""  // Clear the password field after storing
   return nil
}

// Returns user without its hashed password
func (r *userRepository) GetUserByID(ctx context.Context, id int64) (*models.User, error) {
   user := &models.User{}

   query := `SELECT id, username, email FROM user where id = ?`
   err := r.db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Username, &user.Email)
   if err != nil {
      if err == sql.ErrNoRows {
         return nil, fmt.Errorf("user not found")
      }
      return nil, fmt.Errorf("failed to get user: %w", err)
   }

   return user, nil
}

func (r *userRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
   user := &models.User{}

   query := `SELECT id, username, email, password FROM users WHERE username = ?`
   err := r.db.QueryRowContext(ctx, query, username).Scan(&user.ID, &user.Username, &user.Email, &user.Password)
   if err != nil {
      if err == sql.ErrNoRows {
         return nil, fmt.Errorf("user not found")
      }
      return nil, fmt.Errorf("failed to get user: %w", err)
   }

   return user, nil
}

func (r *userRepository) AuthenticateUser(ctx context.Context, username, password string) (*models.User, error) {
   user, err := r.GetUserByUsername(ctx, username)
   if err != nil {
      return nil, err
   }

   err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
   if err != nil {
      return nil, fmt.Errorf("invalid password")
   }

   user.Password = ""  // Clear the password field before returning
   return user, nil
}
