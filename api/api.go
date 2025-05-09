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

type GetLatestUserActivityResult struct {
	Activity []db.GetLatestUserActivityRow
}

//encore:api private method=GET path=/api/user/activity
func GetLatestUserActivity(ctx context.Context, params db.GetLatestUserActivityParams) (*GetLatestUserActivityResult, error) {
	rows, err := db.New().GetLatestUserActivity(ctx, markblogdb.Stdlib(), params)
	if err != nil {
		return nil, err
	}
	res := &GetLatestUserActivityResult{
        Activity: make([]db.GetLatestUserActivityRow, 0),
    }
	for _, r := range rows {
		res.Activity = append(res.Activity, *r)
	}

	return res, nil
}

//encore:api private method=POST path=/api/post
func CreatePost(ctx context.Context, params db.CreatePostParams) (*db.Post, error) {
	return db.New().CreatePost(ctx, markblogdb.Stdlib(), params)
}

//encore:api private method=GET path=/api/post/id/:id
func GetPostByID(ctx context.Context, id uuid.UUID) (*db.Post, error) {
	return db.New().GetPostByID(ctx, markblogdb.Stdlib(), id)
}

type GetLatestPostsResult struct {
	Posts []db.GetLatestPostsRow `json:"posts"`
}

//encore:api private method=GET path=/api/posts/latest
func GetLatestPosts(ctx context.Context, params db.GetLatestPostsParams) (*GetLatestPostsResult, error) {
	rows, err := db.New().GetLatestPosts(ctx, markblogdb.Stdlib(), params)
	if err != nil {
		return nil, err
	}
	posts := &GetLatestPostsResult{
        Posts: make([]db.GetLatestPostsRow, 0),
    }
	for _, r := range rows {
		posts.Posts = append(posts.Posts, *r)
	}

	return posts, nil
}

//encore:api private method=POST path=/api/comment
func CreateComment(ctx context.Context, params db.CreateCommentParams) (*db.Comment, error) {
	return db.New().CreateComment(ctx, markblogdb.Stdlib(), params)
}

//encore:api private method=GET path=/api/comment/id/:id
func GetCommentByID(ctx context.Context, id uuid.UUID) (*db.Comment, error) {
	return db.New().GetCommentByID(ctx, markblogdb.Stdlib(), id)
}

type GetLatestCommentsForPostResult struct {
	Comments []db.GetLatestCommentsForPostRow `json:"comments"`
}

//encore:api private method=GET path=/api/comment
func GetLatestCommentsForPost(ctx context.Context, params db.GetLatestCommentsForPostParams) (*GetLatestCommentsForPostResult, error) {
	rows, err := db.New().GetLatestCommentsForPost(ctx, markblogdb.Stdlib(), params)
	if err != nil {
		return nil, err
	}
	res := &GetLatestCommentsForPostResult{
        Comments: make([]db.GetLatestCommentsForPostRow, 0),
    }
	for _, r := range rows {
		res.Comments = append(res.Comments, *r)
	}

	return res, nil
}