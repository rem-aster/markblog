package app

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"

	db "github.com/rem-aster/markblog/internal/database"
	"github.com/rem-aster/markblog/internal/server"
)

type App struct {
	Env *Env
}

type Env struct {
	SECRET            string
	POSTGRES_USER     string
	POSTGRES_PASSWORD string
	POSTGRES_DB       string
}

func NewApp() *App {
	return &App{}
}

func (a *App) Init() {
	a.SetupEnv()
	a.SetupServer()
}

func (a *App) SetupEnv() {
	a.Env = &Env{
		SECRET:            os.Getenv("SECRET"),
		POSTGRES_USER:     os.Getenv("POSTGRES_USER"),
		POSTGRES_PASSWORD: os.Getenv("POSTGRES_PASSWORD"),
		POSTGRES_DB:       os.Getenv("POSTGRES_DB"),
	}
}

func (a *App) SetupServer() {
	
	ctx := context.Background()

	conn, err := pgx.Connect(ctx, fmt.Sprintf(
		"postgresql://%v:%v@db:5432/%v?sslmode=disable", a.Env.POSTGRES_USER, a.Env.POSTGRES_PASSWORD, a.Env.POSTGRES_DB))
	if err != nil {
		panic(err)
	}
	defer conn.Close(ctx)

	db := db.New(conn)
	h := server.New(db, a.Env.SECRET)
	h.SetupRoutes()

	err = h.StartServer()
	if err != nil {
		log.Printf("Error Starting Server:\nMessage:\n%v", err.Error())
	}
}
