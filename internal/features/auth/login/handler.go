package login

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	svc Service
}

func NewHandler(service Service) Handler {
	return Handler{
		svc: service,
	}
}

func (h *Handler) ServePage(handler string) echo.HandlerFunc {
	page := Page(handler)

	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
		return page.Render(c.Request().Context(), c.Response().Writer)
	}
}

func (h *Handler) HandleFormRequest() echo.HandlerFunc {
	return func(c echo.Context) error {

		login := c.FormValue("login")
		password := c.FormValue("password")

		log.Println("got: ", login, " and ", password)

		_, err := h.svc.Login(c.Request().Context(), login, password)
		if err != nil {
			return c.HTML(http.StatusOK, "<h2>Provided credentials are invalid</h2>")
		}

		c.Response().Header().Add("HX-Redirect", "/")
		return c.HTML(http.StatusOK, "")
	}
}
