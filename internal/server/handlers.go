package server

import (
	"net/http"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"

	db "github.com/rem-aster/markblog/internal/database"
)

func (h *Handler) HealthCheck(e echo.Context) error {
	return e.JSON(http.StatusOK, "OK")
}

func (h *Handler) FormLogin(e echo.Context) error {
    var credentials struct {
        Username string `json:"username"`
        Password string `json:"password"`
    }
    if err := e.Bind(&credentials); err != nil {
        credentials.Username = e.FormValue("username")
        credentials.Password = e.FormValue("password")
    }
    if r, err := h.db.CheckUserExists(e.Request().Context(), credentials.Username); (err != nil) || !r {
        return e.JSON(http.StatusUnauthorized, map[string]string{"error": "Unknown username"})
    }

    user, err := h.db.GetuserByUsername(e.Request().Context(), credentials.Username)
    if err != nil {
        return err
    }
    if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(credentials.Password)); err != nil {
        return e.JSON(http.StatusUnauthorized, map[string]string{"error": "Wrong password"})
    }
    sess, _ := session.Get("session", e)
    sess.Values["authenticated"] = true
    sess.Values["username"] = credentials.Username
    if err := sess.Save(e.Request(), e.Response()); err != nil {
        return e.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save session"})
    }
    return e.JSON(http.StatusOK, map[string]string{"message": "Login successful", "username": credentials.Username})
}

func (h *Handler) FormRegister(e echo.Context) error {
    var credentials struct {
        Username string `json:"username" form:"username"`
        Password string `json:"password" form:"password"`
    }
    if err := e.Bind(&credentials); err != nil {
        return e.JSON(http.StatusBadRequest, map[string]string{
            "error": "Invalid request format",
        })
    }
    if credentials.Username == "" || credentials.Password == "" {
        return e.JSON(http.StatusBadRequest, map[string]string{
            "error": "Username and password are required",
        })
    }
    if len(credentials.Password) < 8 {
        return e.JSON(http.StatusBadRequest, map[string]string{
            "error": "Password must be at least 8 characters",
        })
    }
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(credentials.Password), bcrypt.DefaultCost)
    if err != nil {
        return e.JSON(http.StatusInternalServerError, map[string]string{
            "error": "Failed to secure password",
        })
    }
    exists, err := h.db.CheckUserExists(e.Request().Context(), credentials.Username)
    if err != nil {
        return e.JSON(http.StatusInternalServerError, map[string]string{
            "error": "Database error",
        })
    }
    if exists {
        return e.JSON(http.StatusConflict, map[string]string{
            "error": "Username already taken",
        })
    }
    _, err = h.db.CreateUser(e.Request().Context(), db.CreateUserParams{
        Username:     credentials.Username,
        PasswordHash: string(hashedPassword),
    })
    if err != nil {
        return e.JSON(http.StatusInternalServerError, map[string]string{
            "error": "Failed to create user",
        })
    }
    sess, err := session.Get("session", e)
    if err != nil {
        return e.JSON(http.StatusInternalServerError, map[string]string{
            "error": "Failed to initialize session",
        })
    }
    sess.Values["authenticated"] = true
    sess.Values["username"] = credentials.Username
    if err := sess.Save(e.Request(), e.Response()); err != nil {
        return e.JSON(http.StatusInternalServerError, map[string]string{
            "error": "Failed to save session",
        })
    }
    return e.JSON(http.StatusCreated, map[string]string{
        "message":  "Registration successful",
        "username": credentials.Username,
    })
}