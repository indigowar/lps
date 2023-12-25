package dashboard

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	db *sqlx.DB
}

func NewHandler(db *sqlx.DB) Handler {
	return Handler{
		db: db,
	}
}

func (h *Handler) ShowDashboard() echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
		return Page().Render(c.Request().Context(), c.Response().Writer)
	}
}
