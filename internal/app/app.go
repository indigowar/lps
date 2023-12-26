package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	session "github.com/spazzymoto/echo-scs-session"

	"lps/internal/config"
	"lps/internal/features/auth"
	"lps/internal/features/dashboard"
	"lps/internal/features/profile"
	"lps/pkg/postgres"
)

func Run(cfg *config.Config) {
	var p *sqlx.DB
	var err error
	if os.Getenv("POSTGRES_URL") != "" {
		p, err = postgres.CreateConnectionUsingURL(os.Getenv("POSTGRES_URL"))
	} else {
		p, err = postgres.CreateConnection(cfg.Db.Host, cfg.Db.Port, cfg.Db.Db, cfg.Db.SystemUser, cfg.Db.SystemPassword)
	}
	if err != nil {
		log.Fatal(err)
	}

	sessionManager := scs.New()
	sessionManager.Lifetime = 24 * time.Hour

	e := echo.New()
	e.Logger.SetLevel(log.INFO)

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"*"},
	}))
	e.Use(session.LoadAndSave(sessionManager))

	e.GET("/status", func(c echo.Context) error {
		return c.HTML(http.StatusOK, "<h1>OK</h1>")
	})

	authService := auth.NewPostgresService(p)
	authHandler := auth.NewHandler(authService, sessionManager)
	e.GET("/auth/login", authHandler.ServeLoginPage("/auth/login"))
	e.POST("/auth/login", authHandler.HandleLoginRequest())
	e.GET("/auth/register/:login", authHandler.ServeRegisterPage("/auth/register"))
	e.POST("/auth/register", authHandler.HandleRegisterRequest())
	e.GET("/auth/logout", authHandler.HandleLogout())

	profileService := profile.NewPostgresService(p)
	profileHandler := profile.NewHandler(profileService, sessionManager)
	e.GET("/profile", profileHandler.GetProfile("/auth/login", "/profile/account", "/profile/employee"))
	e.GET("/profile/account", profileHandler.ServeAccountUpdateForm("/profile/account", "/profile/cancel"))
	e.PUT("/profile/account", profileHandler.HandleAccountUpdate("/profile"))
	e.GET("/profile/employee", profileHandler.ServeEmployeeUpdateForm("/profile/employee", "/profile/cancel"))
	e.PUT("/profile/employee", profileHandler.HandleEmployeeUpdate("/profile"))
	e.GET("/profile/cancel", profileHandler.HandleEditCancellation("/profile"))

	dashboardService := dashboard.NewPostgrseService(p)
	dashboardHandler := dashboard.NewHandler(dashboardService, sessionManager)
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
