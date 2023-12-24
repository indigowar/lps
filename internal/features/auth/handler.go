package auth

import (
	"log"
	"net/http"

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
		// currentSession, _ := session.Get("user-session", c)
		// id, exists := currentSession.Values["user-id"]
		// log.Println(id.(string))
		// if exists {
		// 	return c.Redirect(http.StatusPermanentRedirect, "/")
		// }

		c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
		return loginPage(handler).Render(c.Request().Context(), c.Response().Writer)
	}
}

func (h *Handler) HandleLoginRequest() echo.HandlerFunc {
	return func(c echo.Context) error {
		// currentSession, _ := session.Get("user-session", c)
		// if _, exists := currentSession.Values["user-id"]; exists {
		// 	log.Println("got logged in user")
		// 	return c.Redirect(http.StatusPermanentRedirect, "/")
		// }

		login := c.FormValue("login")
		password := c.FormValue("password")

		_, err := h.svc.Login(c.Request().Context(), login, password)
		if err != nil {
			log.Printf("Failed to login, due to: %s\n", err)
			return echo.NewHTTPError(http.StatusBadRequest, "credentials are invalid")
		}

		// currentSession.Values["user-id"] = id.String()
		// _ = currentSession.Save(c.Request(), c.Response().Writer)

		c.Response().Header().Add("HX-Redirect", "/")
		return c.NoContent(http.StatusOK)
	}
}

func (h *Handler) ServeRegisterPage(handler string) echo.HandlerFunc {
	return func(c echo.Context) error {
		// currentSession, _ := session.Get("user-session", c)
		// _, exists := currentSession.Values["user-id"]
		// if exists {
		// 	return c.Redirect(http.StatusPermanentRedirect, "/")
		// }

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

func (h *Handler) HandleLogout() echo.HandlerFunc {
	return func(c echo.Context) error {
		// currentSession, _ := session.Get("user-session", c)

		// if _, exists := currentSession.Values["user-id"]; !exists {
		// 	return echo.NewHTTPError(http.StatusUnauthorized, "<h1>unauthorized</h1>")
		// }
		// delete(currentSession.Values, "user-id")
		// currentSession.Save(c.Request(), c.Response().Writer)
		return c.NoContent(http.StatusOK)
	}
}
