package db

import (
	"database/sql"
	"errors"
	"fmt"
	"marketplace/src/database/repository"
	"time"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

func Connect(conf repository.Config) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		conf.Host,
		conf.Port,
		conf.User,
		conf.Password,
		conf.Name,
	)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}
	return db, nil
}

func Close(db *sql.DB) (bool, error) {
	err := db.Close()
	return err == nil, err
}
func Ping(db *sql.DB) (bool, error) {
	ping := db.Ping()
	return ping == nil, ping
}

type PostgresUserRepo struct {
	db *sql.DB
}

func UserRepo(db *sql.DB) *PostgresUserRepo {
	return &PostgresUserRepo{db: db}
}

func (r *PostgresUserRepo) Create(user *repository.User) error {
	user.ID = fmt.Sprintf("user_%d", time.Now().UnixNano())

	query := `INSERT INTO users (id, username, email, password_hash) VALUES ($1, $2, $3, $4) RETURNING created_at`
	err := r.db.QueryRow(query, user.ID, user.Username, user.Email, user.PasswordHash).Scan(&user.CreatedAt)
	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			switch pgErr.Constraint {
			case "users_email_key":
				return repository.ErrDuplicateEmail
			case "users_username_key":
				return repository.ErrDuplicateUsername
			}
		}
		return err
	}
	return nil
}

func (r *PostgresUserRepo) GetByID(id string) (*repository.User, error) {
	user := &repository.User{}
	query := `SELECT id, username, email, password_hash, created_at FROM users WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}
