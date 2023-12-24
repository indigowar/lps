package auth

import (
	"log"
	"net/http"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) ServeLoginPage(handler string) echo.HandlerFunc {
	return func(c echo.Context) error {
		currentSession, _ := session.Get("user-session", c)
		_, exists := currentSession.Values["user-id"]
		if exists {
			return c.Redirect(http.StatusPermanentRedirect, "/")
		}

		c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
		return loginPage(handler).Render(c.Request().Context(), c.Response().Writer)
	}
}

func (h *Handler) HandleLoginRequest() echo.HandlerFunc {
	return func(c echo.Context) error {
		currentSession, _ := session.Get("user-session", c)
		_, exists := currentSession.Values["user-id"]
		if exists {
			return c.Redirect(http.StatusPermanentRedirect, "/")
		}

		login := c.FormValue("login")
		password := c.FormValue("password")

		log.Println("got: ", login, " and ", password)

		id, err := h.svc.Login(c.Request().Context(), login, password)
		if err != nil {
			return c.HTML(http.StatusOK, "<h2>Provided credentials are invalid</h2>")
		}

		currentSession.Values["user-id"] = id.String()
		if err = currentSession.Save(c.Request(), c.Response().Writer); err != nil {
			log.Println("FAILED TO SAVE THE FREAKING SESSION: ", err)
		}

		log.Println("CREATED AND EVERYTHING SHOULD FUCKING WORK")

		c.Response().Header().Add("HX-Redirect", "/")
		return c.HTML(http.StatusOK, "")
	}
}

func (h *Handler) ServeRegisterPage(handler string) echo.HandlerFunc {
	return func(c echo.Context) error {
		currentSession, _ := session.Get("user-session", c)
		_, exists := currentSession.Values["user-id"]
		if exists {
			return c.Redirect(http.StatusPermanentRedirect, "/")
		}

		login := c.Param("login")

		c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
		return registrationPage(handler, login).Render(c.Request().Context(), c.Response().Writer)
	}
}

func (h *Handler) HandleRegisterRequest() echo.HandlerFunc {
	return func(c echo.Context) error {
		login := c.FormValue("login")
		password := c.FormValue("password")

		log.Println("Found: ", login, " and ", password)

		return nil
	}
}
