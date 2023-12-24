package dashboard

import (
	"github.com/labstack/echo/v4"
)

type Handler struct{}

func NewHandler() Handler {
	return Handler{}
}

func (h *Handler) ShowDashboard() echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
		return Page().Render(c.Request().Context(), c.Response().Writer)
	}
}
