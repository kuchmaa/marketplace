package repository

import (
	"errors"
	"time"
)

type UserService interface {
	Create(user *User) error
	GetByID(id string) (*User, error)
}

type User struct {
	ID           string
	Username     string
	Email        string
	PasswordHash string
	CreatedAt    time.Time
}

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

var (
	ErrDuplicateEmail    = errors.New("email already exists")
	ErrDuplicateUsername = errors.New("username already exists")
)
