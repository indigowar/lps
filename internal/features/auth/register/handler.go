package register

import (
	"log"

	"github.com/labstack/echo/v4"
)

type Handler struct{}

func NewHandler() Handler {
	return Handler{}
}

func (h *Handler) ServePage(handler string) echo.HandlerFunc {
	page := Page(handler, "Nagibator2008")

	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
		return page.Render(c.Request().Context(), c.Response().Writer)
	}
}

func (h *Handler) HandleFormRequest() echo.HandlerFunc {
	return func(c echo.Context) error {
		login := c.FormValue("login")
		password := c.FormValue("password")

		log.Println("Found: ", login, " and ", password)

		return nil
	}
}
