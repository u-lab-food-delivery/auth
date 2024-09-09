package postgres

import (
	"auth_service/models"
	"auth_service/storage/cache"
	"context"
	"database/sql"
	"fmt"
	"log"

	sq "github.com/Masterminds/squirrel"
	"golang.org/x/crypto/bcrypt"
)

type UserManagementImpl struct {
	db         *sql.DB
	sqlBuilder sq.StatementBuilderType
	cache      *cache.AuthCache
}

func NewUserManagementSQL(db *sql.DB, cache *cache.AuthCache) *UserManagementImpl {
	return &UserManagementImpl{
		db:         db,
		sqlBuilder: sq.StatementBuilderType{}.PlaceholderFormat(sq.Dollar),
		cache:      cache,
	}
}

func (a *UserManagementImpl) CreateUser(ctx context.Context, req *models.User) (*models.User, error) {
	// Hash password and insert new user into the database
	hashedPassword, err := HashPassword(req.HashedPassword)
	if err != nil {
		return nil, err
	}

	err = a.InsertUser(ctx, req.Email, hashedPassword, req.Name)
	if err != nil {
		return nil, err
	}

	// Retrieve the newly created user
	user, err := a.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	// Cache the new user by both ID and email
	err = a.cache.CreateOrUpdateUserByID(ctx, user)
	if err != nil {
		log.Println("Redis Error: ", err)
		return nil, err
	}

	err = a.cache.CreateOrUpdateUserByEmail(ctx, user)
	if err != nil {
		log.Println("Redis Error: ", err)
		return nil, err
	}

	return user, nil
}

func (a *UserManagementImpl) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	// Check cache by email first
	cachedUser, err := a.cache.GetUserByEmail(ctx, email)
	if err != nil {
		log.Println("Redis Error: ", err)
		return nil, err
	}

	if cachedUser != nil {
		return cachedUser, nil
	}

	// Check database
	sqlQuery, args, err := a.sqlBuilder.Select(
		"user_id",
		"email",
		"hashed_password",
		"name",
		"created_at",
		"updated_at",
	).From("users").Where(
		sq.Eq{"email": email},
	).ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build SQL query: %v", err)
	}

	row := a.db.QueryRowContext(ctx, sqlQuery, args...)
	user := &models.User{}
	err = row.Scan(
		&user.UserId,
		&user.Email,
		&user.HashedPassword,
		&user.Name,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		log.Println("Failed to scan user: ", err)
		return nil, err
	}

	// Cache the retrieved user by both ID and email
	err = a.cache.CreateOrUpdateUserByEmail(ctx, user)
	if err != nil {
		log.Println("Redis Error: ", err)
		return nil, err
	}

	err = a.cache.CreateOrUpdateUserByID(ctx, user)
	if err != nil {
		log.Println("Redis Error: ", err)
		return nil, err
	}

	return user, nil
}

func (a *UserManagementImpl) UpdateUser(ctx context.Context, user *models.User) (*models.User, error) {
	// Update database
	sqlQuery, args, err := a.sqlBuilder.Update("users").
		Set("email", user.Email).
		Set("hashed_password", user.HashedPassword).
		Set("name", user.Name).
		Set("updated_at", user.UpdatedAt).
		Where(sq.Eq{"user_id": user.UserId}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build SQL query: %v", err)
	}

	_, err = a.db.ExecContext(ctx, sqlQuery, args...)
	if err != nil {
		log.Println("Failed to update user: ", err)
		return nil, err
	}

	// Update cache by both ID and email
	err = a.cache.CreateOrUpdateUserByID(ctx, user)
	if err != nil {
		log.Println("Redis Error: ", err)
		return nil, err
	}

	err = a.cache.CreateOrUpdateUserByEmail(ctx, user)
	if err != nil {
		log.Println("Redis Error: ", err)
		return nil, err
	}

	return user, nil
}

func (a *UserManagementImpl) DeleteUser(ctx context.Context, userId string) error {
	// Get user by ID to retrieve email for cache deletion
	user, err := a.GetByID(ctx, userId)
	if err != nil {
		return err
	}

	// Delete from database
	sqlQuery, args, err := a.sqlBuilder.Delete("users").
		Where(sq.Eq{"user_id": userId}).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build SQL query: %v", err)
	}

	_, err = a.db.ExecContext(ctx, sqlQuery, args...)
	if err != nil {
		log.Println("Failed to delete user: ", err)
		return err
	}

	// Delete from cache by both ID and email
	err = a.cache.DeleteUserByID(ctx, userId)
	if err != nil {
		log.Println("Redis Error: ", err)
		return err
	}

	err = a.cache.DeleteUserByEmail(ctx, user.Email)
	if err != nil {
		log.Println("Redis Error: ", err)
		return err
	}

	return nil
}

func (a *UserManagementImpl) InsertUser(ctx context.Context, email, hashedPassword, name string) error {
	sqlQuery, args, err := a.sqlBuilder.Insert("users").
		Columns("email", "hashed_password", "name", "created_at").
		Values(email, hashedPassword, name, "NOW()").
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build SQL query: %v", err)
	}

	_, err = a.db.ExecContext(ctx, sqlQuery, args...)
	if err != nil {
		log.Println("Failed to insert user: ", err)
		return err
	}

	return nil
}

func (a *UserManagementImpl) GetByID(ctx context.Context, userId string) (*models.User, error) {
	// Check cache by ID first
	cachedUser, err := a.cache.GetUserByID(ctx, userId)
	if err != nil {
		log.Println("Redis Error: ", err)
		return nil, err
	}

	if cachedUser != nil {
		return cachedUser, nil
	}

	// Check database
	sqlQuery, args, err := a.sqlBuilder.Select(
		"user_id",
		"email",
		"hashed_password",
		"name",
		"created_at",
		"updated_at",
	).From("users").Where(
		sq.Eq{"user_id": userId},
	).ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build SQL query: %v", err)
	}

	row := a.db.QueryRowContext(ctx, sqlQuery, args...)
	user := &models.User{}
	err = row.Scan(
		&user.UserId,
		&user.Email,
		&user.HashedPassword,
		&user.Name,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		log.Println("Failed to scan user: ", err)
		return nil, err
	}

	// Cache the retrieved user by both ID and email
	err = a.cache.CreateOrUpdateUserByEmail(ctx, user)
	if err != nil {
		log.Println("Redis Error: ", err)
		return nil, err
	}

	err = a.cache.CreateOrUpdateUserByID(ctx, user)
	if err != nil {
		log.Println("Redis Error: ", err)
		return nil, err
	}

	return user, nil
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error hashing password:", err)
		return "", err
	}
	return string(hashedPassword), nil
}

func ComparePassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return false
		}
		log.Println("Error comparing passwords:", err)
		return false
	}
	return true
}
