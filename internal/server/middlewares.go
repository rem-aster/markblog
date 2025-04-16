package server

import (
	"net/http"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func AuthMiddleware() echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(e echo.Context) error {
            sess, err := session.Get("session", e)
            if err != nil {
                return e.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
            }
            if auth, ok := sess.Values["authenticated"].(bool); !ok || !auth {
                return e.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
            }
            return next(e)
        }
    }
}