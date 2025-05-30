// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0

package db

import (
	"time"

	"encore.dev/types/uuid"
)

type Comment struct {
	ID        uuid.UUID
	PostID    uuid.UUID
	UserID    *uuid.UUID
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Post struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type User struct {
	ID           uuid.UUID
	Username     string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
