package server

import (
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	db "github.com/rem-aster/markblog/internal/database"
)

type Handler struct {
	srv    *echo.Echo
	db     *db.Queries
	secret string
}

func New(db *db.Queries, secret string) *Handler {
	return &Handler{
		srv:    echo.New(),
		db:     db,
		secret: secret,
	}
}

func (h *Handler) SetupRoutes() {
	store := sessions.NewCookieStore([]byte(h.secret))
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
		// Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}

	h.srv.Use(
		middleware.Logger(),
		middleware.Recover(),
		session.Middleware(store),
	)

	h.srv.Static("/", "web/dist").Name = "static"
	h.srv.File("/", "web/dist/index.html").Name = "index"
	h.srv.GET("/health", h.HealthCheck).Name = "healthcheck"

	l := h.srv.Group("", AuthMiddleware())
	l.GET("/lol", h.HealthCheck)
}

func (h *Handler) StartServer() error {
	h.srv.Debug = true
	return h.srv.Start(":8080")
}