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
	"lps/pkg/postgres"
)

func Run(cfg *config.Config) {
	_, err := postgres.CreateClient(postgres.Credentials{
		Host:     cfg.Db.Host,
		Port:     cfg.Db.Port,
		User:     cfg.Db.SystemUser,
		Password: cfg.Db.SystemPassword,
	})

	if err != nil {
		log.Fatal(err)
	}

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

	_ = dashboard.NewHandler()
	// e.GET("/", dashboardHandler.ShowDashboard())

	profileHandler := profile.NewHandler()
	e.GET("/profile", profileHandler.GetProfile("/auth/login"))
	e.PUT("/profile", profileHandler.ServeProfileRequst())

	e.GET("/", func(c echo.Context) error {
		sess, _ := session.Get("user-sesson", c)
		_, exists := sess.Values["user-id"]
		if !exists {
			return c.HTML(http.StatusOK, "Not Logged In")
		}
		return c.HTML(http.StatusOK, "Logged In")
	})

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