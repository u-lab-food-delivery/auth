package postgres

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/go-redis/redis/v8"
)

type AuthManagementImpl struct {
	db         *sql.DB
	sqlBuilder sq.StatementBuilderType
	redis      *redis.Client
}

func NewAuthManagementSQL(db *sql.DB, redis *redis.Client) *AuthManagementImpl {
	return &AuthManagementImpl{
		db:         db,
		sqlBuilder: sq.StatementBuilderType{}.PlaceholderFormat(sq.Dollar),
		redis:      redis,
	}
}

func (a *AuthManagementImpl) Register()
