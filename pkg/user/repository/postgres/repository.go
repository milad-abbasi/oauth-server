package postgres

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/milad-abbasi/oauth-server/pkg/user"
	"go.uber.org/zap"
)

const usersTable = "users"

var pgBuilder = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
var pgDefault = sq.Expr("DEFAULT")

type Repository struct {
	logger *zap.Logger
	pgPool *pgxpool.Pool
}

func New(logger *zap.Logger, pgPool *pgxpool.Pool) *Repository {
	return &Repository{
		logger: logger.Named("UserRepository"),
		pgPool: pgPool,
	}
}

func (r *Repository) Exists(ctx context.Context, email string) (bool, error) {
	query, args, err := pgBuilder.
		Select("TRUE").
		From(usersTable).
		Where(sq.Eq{"email": email}).
		ToSql()
	if err != nil {
		return false, fmt.Errorf("invalid query: %w", err)
	}

	var exists bool
	err = r.pgPool.
		QueryRow(ctx, query, args...).
		Scan(&exists)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	} else if err != nil {
		return false, fmt.Errorf("query failed: %w", err)
	}

	return exists, nil
}

func (r *Repository) FindOne(ctx context.Context, u *user.User) (*user.User, error) {
	qBuilder := pgBuilder.Select("*").From(usersTable)
	if u.ID != "" {
		qBuilder = qBuilder.Where(sq.Eq{"id": u.ID})
	}
	if u.Name != "" {
		qBuilder = qBuilder.Where(sq.Like{"name": u.Name})
	}
	if u.Email != "" {
		qBuilder = qBuilder.Where(sq.Eq{"email": u.Email})
	}

	query, args, err := qBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("invalid query: %w", err)
	}

	err = r.pgPool.
		QueryRow(ctx, query, args...).
		Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.CreatedAt, &u.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, user.ErrUserNotFound
	} else if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	return u, nil
}

func (r *Repository) Create(ctx context.Context, u *user.User) (*user.User, error) {
	query, args, err := pgBuilder.
		Insert(usersTable).
		Columns("name", "email", "password").
		Values(u.Name, u.Email, u.Password).
		Suffix("RETURNING id, created_at, updated_at").
		ToSql()
	if err != nil {
		return nil, err
	}

	err = r.pgPool.
		QueryRow(ctx, query, args...).
		Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt)

	if err != nil {
		var pgErr *pgconn.PgError
		if ok := errors.As(err, &pgErr); ok && pgErr.Code == pgerrcode.UniqueViolation {
			return nil, user.ErrUserExists
		}

		return nil, err
	}

	return u, nil
}
