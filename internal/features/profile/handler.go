package profile

import (
	"net/http"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

type Handler struct {
}

func NewHandler() Handler {
	return Handler{}
}

func (h *Handler) GetProfile(loginPage string) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, _ := session.Get("user-session", c)

		_, exists := sess.Values["user-id"]
		if !exists {
			return c.Redirect(http.StatusPermanentRedirect, loginPage)
		}

		c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
		return Page().Render(c.Request().Context(), c.Response().Writer)
	}
}

func (h *Handler) ServeProfileRequst() echo.HandlerFunc {
	return func(c echo.Context) error {
		return nil
	}
}
