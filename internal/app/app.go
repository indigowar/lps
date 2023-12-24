package app

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"

	"lps/internal/config"
	"lps/internal/features/auth/login"
	"lps/internal/features/auth/register"
	"lps/internal/features/dashboard"
)

func Run(cfg *config.Config) {
	e := echo.New()
	e.Logger.SetLevel(log.INFO)

	e.Use(session.Middleware(sessions.NewCookieStore([]byte("secret"))))

	e.GET("/status", func(c echo.Context) error {
		return c.HTML(http.StatusOK, "<h1>OK</h1>")
	})

	loginService := loginService{
		Users: []user{
			{uuid.New(), "nagibator44", "123456"},
			{uuid.New(), "nagibator45", "abcd"},
		},
	}

	loginHandler := login.NewHandler(&loginService)
	e.GET("/auth/login", loginHandler.ServePage("/auth/login"))
	e.POST("/auth/login", loginHandler.HandleFormRequest())

	registerHandler := register.NewHandler()
	e.GET("/auth/register", registerHandler.ServePage("/auth/register"))
	e.POST("/auth/register", registerHandler.HandleFormRequest())

	dashboardHandler := dashboard.NewHandler()
	e.GET("/", dashboardHandler.ShowDashboard())

	go func() {
		if err := e.Start(":3000"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("Shutting down the server, because: %s", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}

type user struct {
	ID       uuid.UUID
	Login    string
	Password string
}

type loginService struct {
	Users []user
}

func (svc *loginService) Login(_ context.Context, login string, password string) (uuid.UUID, error) {
	for _, u := range svc.Users {
		if u.Login == login && u.Password == password {
			return u.ID, nil
		}
	}
	return uuid.UUID{}, errors.New("credentials are invalid")
}
