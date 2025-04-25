package api

import (
	"context"

	"encore.app/api/db"
	"encore.dev/storage/sqldb"
	"encore.dev/types/uuid"
)

var markblogdb = sqldb.NewDatabase("markblog", sqldb.DatabaseConfig{
	Migrations: "./db/migrations",
})

//encore:api private method=POST path=/api/user
func CreateUser(ctx context.Context, params db.CreateUserParams) (*db.User, error) {
	return db.New().CreateUser(ctx, markblogdb.Stdlib(), params)
}

//encore:api private method=GET path=/api/user/id/:id
func GetUserByID(ctx context.Context, id uuid.UUID) (*db.User, error) {
	return db.New().GetUserByID(ctx, markblogdb.Stdlib(), id)
}

//encore:api private method=GET path=/api/user/username/:username
func GetUserByUsername(ctx context.Context, username string) (*db.User, error) {
	return db.New().GetUserByUsername(ctx, markblogdb.Stdlib(), username)
}

type CheckUserExistsResult struct {
	Exists bool `json:"exists"`
}

//encore:api private method=GET path=/api/user/exists/:username
func CheckUserExists(ctx context.Context, username string) (*CheckUserExistsResult, error) {
	res := new(CheckUserExistsResult)
	var err error
	res.Exists, err = db.New().CheckUserExists(ctx, markblogdb.Stdlib(), username)
	return res, err
}
