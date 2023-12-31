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
	access_error "lps/internal/features/access_error"
	addincident "lps/internal/features/add_incident"
	addworker "lps/internal/features/add_worker"
	"lps/internal/features/auth"
	"lps/internal/features/dashboard"
	"lps/internal/features/dashboard/departments"
	staff "lps/internal/features/dashboard/employee"
	"lps/internal/features/dashboard/incidents"
	"lps/internal/features/dashboard/positions"
	"lps/internal/features/profile"
	postgresUsecases "lps/internal/repository/postgres"
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
	e.GET("/position", dashboardHandler.ServePositionsTable())

	departmentUseCase := postgresUsecases.NewDepartmentUsecase(p)
	departmentHandler := departments.NewHandler(departmentUseCase)
	e.GET("/department", departmentHandler.ServeTable())
	e.GET("/department/:id", departmentHandler.ServeItem())
	e.GET("/department/:id/edit", departmentHandler.ServeEdit())
	e.PUT("/department", departmentHandler.HandleEdit())
	e.DELETE("/department/:id", departmentHandler.Delete())

	positionsUseCase := postgresUsecases.NewPositionUseCase(p)
	positionHandler := positions.NewHandler(positionsUseCase)
	e.GET("/position", positionHandler.ServeTable())
	e.GET("/position/:id", positionHandler.ServeItem())
	e.GET("/position/:id/edit", positionHandler.ServeEdit())
	e.PUT("/position", positionHandler.HandleEdit())
	e.DELETE("/position/:id", positionHandler.Delete())

	staffUseCase := postgresUsecases.NewEmployeeUseCase(p)
	staffHandler := staff.NewHandler(staffUseCase, departmentUseCase, positionsUseCase)
	e.GET("/staff", staffHandler.ServeTable())
	e.GET("/staff/:id", staffHandler.ServeItem())
	e.GET("/staff/:id/edit", staffHandler.ServeEdit())
	e.PUT("/staff", staffHandler.HandleEdit())
	e.DELETE("/staff/:id", staffHandler.Delete())

	addWorkerService := addworker.NewPostgresService(p)
	addWorkerHandler := addworker.NewHandler(addWorkerService, sessionManager)
	e.GET("/add-worker", addWorkerHandler.ServePage("/add-worker"))
	e.POST("/add-worker", addWorkerHandler.HandleRequest())

	createIncidentUseCase := postgresUsecases.NewCreateIncidentsUseCase(p)
	addIncidentHandler := addincident.NewHandler(sessionManager, createIncidentUseCase)
	e.GET("/add-incident", addIncidentHandler.ServePage())
	e.POST("/add-incident", addIncidentHandler.HandleRequest())

	incidentUseCase := postgresUsecases.NewIncidentUseCase(p)
	incidentHandler := incidents.NewIncidentHandler(incidentUseCase)

	e.GET("/incident", incidentHandler.ServeTable())
	e.PUT("/incident/:id", incidentHandler.HandleEdit())
	e.DELETE("/incident/:id", incidentHandler.Delete())
	e.GET("/incident/:id/edit", incidentHandler.ServeEdit())

	e.GET("/denied", func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
		return access_error.AccessError().Render(c.Request().Context(), c.Response().Writer)
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
