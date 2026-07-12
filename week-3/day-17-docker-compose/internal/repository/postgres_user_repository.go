package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/meedaycodes/day14-integration-testing/internal/model"
)

// PostgresUserRepository implements UserRepository using PostgreSQL as the
// storage backend. It holds a connection pool (*pgxpool.Pool), not a single
// connection — the pool manages multiple connections and is safe for concurrent use.
type PostgresUserRepository struct {
	pool *pgxpool.Pool
}

// NewPostgresUserRepository creates a new PostgresUserRepository with the given
// connection pool. The pool is created externally (in main.go) and injected here —
// the repository uses the pool but doesn't own or configure it.
func NewPostgresUserRepository(pool *pgxpool.Pool) *PostgresUserRepository {
	newPostgresUserRepository := &PostgresUserRepository{pool: pool}
	return newPostgresUserRepository
}

// Save inserts a new user into the users table. Uses parameterized queries ($1, $2, $3)
// to prevent SQL injection. Exec is used because INSERT doesn't return rows.
func (p *PostgresUserRepository) Save(ctx context.Context, user model.User) error {

	sqlStatement := "INSERT INTO users (id, name, email, password_hash) VALUES ($1,$2,$3,$4)"

	cmdTag, err := p.pool.Exec(ctx, sqlStatement, user.ID, user.Name, user.Email, user.PasswordHash)

	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("user with id %s was not saved in the users table", user.ID)

	}
	return nil
}

// FindByID retrieves a single user by ID using QueryRow (expects one row).
// Scan reads the row columns into the User struct fields. If no row matches,
// pgx returns pgx.ErrNoRows which we translate to ErrUserNotFound to keep
// error handling consistent across repository implementations.
func (p *PostgresUserRepository) FindByID(ctx context.Context, id string) (model.User, error) {

	var user model.User

	sqlStatement := "SELECT id, name, email,password_hash  FROM users WHERE id = $1"

	row := p.pool.QueryRow(ctx, sqlStatement, id)
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash)

	if errors.Is(err, pgx.ErrNoRows) {
		return user, ErrUserNotFound
	} else if err != nil {
		return user, err
	}

	return user, nil
}

// FindAll retrieves all users from the database. Uses Query (not QueryRow)
// because it returns multiple rows. rows.Close is deferred to release the
// database connection back to the pool. rows.Err() is checked after the loop
// to catch errors that occurred during iteration.
func (p *PostgresUserRepository) FindAll(ctx context.Context, limit, offset int) ([]model.User, error) {

	var users []model.User

	sqlStatement := "SELECT id, name, email, password_hash FROM users LIMIT $1 OFFSET $2"
	rows, err := p.pool.Query(ctx, sqlStatement, limit, offset)

	if err != nil {
		return users, err
	}

	defer rows.Close()

	for rows.Next() {

		var user model.User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash)

		if err != nil {
			return users, err
		}
		users = append(users, user)

	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating user rows: %w", err)
	}

	return users, nil

}

// Update modifies an existing user's name and email. Uses RowsAffected to detect
// whether the ID existed — if zero rows were affected, the user wasn't found.
func (p *PostgresUserRepository) Update(ctx context.Context, user model.User) error {

	sqlStatement := "UPDATE users SET name = $1 , email = $2 WHERE id = $3"
	cmdTag, err := p.pool.Exec(ctx, sqlStatement, user.Name, user.Email, user.ID)

	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return ErrUserNotFound
	}

	return nil

}

// Delete removes a user by ID. Like Update, checks RowsAffected to determine
// if the user existed. Returns ErrUserNotFound if no rows were deleted.
func (p *PostgresUserRepository) Delete(ctx context.Context, id string) error {

	sqlStatement := "DELETE FROM users WHERE id = $1"

	cmdTag, err := p.pool.Exec(ctx, sqlStatement, id)

	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return ErrUserNotFound
	}

	return nil
}

// FindByEmail retrieves a single user by Email using QueryRow (expects one row).
// Scan reads the row columns into the User struct fields. If no row matches,
// pgx returns pgx.ErrNoRows which we translate to ErrUserNotFound to keep
// error handling consistent across repository implementations.
func (p *PostgresUserRepository) FindByEmail(ctx context.Context, email string) (model.User, error) {

	var user model.User

	sqlStatement := "SELECT id, name, email,password_hash  FROM users WHERE email = $1"

	row := p.pool.QueryRow(ctx, sqlStatement, email)
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash)

	if errors.Is(err, pgx.ErrNoRows) {
		return user, ErrUserNotFound
	} else if err != nil {
		return user, err
	}

	return user, nil
}
