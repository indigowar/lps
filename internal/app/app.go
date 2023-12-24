package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"

	"lps/internal/config"
	"lps/internal/features/auth"
	"lps/internal/features/dashboard"
	"lps/internal/features/profile"
)

func Run(cfg *config.Config) {
	e := echo.New()
	e.Logger.SetLevel(log.INFO)

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"*"},
	}))
	e.Use(session.Middleware(sessions.NewCookieStore([]byte("secret"))))

	e.GET("/status", func(c echo.Context) error {
		return c.HTML(http.StatusOK, "<h1>OK</h1>")
	})

	authHandler := auth.NewHandler(nil)
	e.GET("/auth/login", authHandler.ServeLoginPage("/auth/login"))
	e.POST("/auth/login", authHandler.HandleLoginRequest())
	e.GET("/auth/register/:login", authHandler.ServeRegisterPage("/auth/register"))
	e.POST("/auth/register", authHandler.HandleRegisterRequest())

	dashboardHandler := dashboard.NewHandler()
	e.GET("/", dashboardHandler.ShowDashboard())

	profileHandler := profile.NewHandler()
	e.GET("/profile", profileHandler.GetProfile("/auth/login"))
	e.PUT("/profile", profileHandler.ServeProfileRequst())

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
